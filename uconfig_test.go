package uconfig_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins"
)

type BadPlugin interface {
	plugins.Plugin

	NotWalkerOrVisitor()
}

func TestBadPlug(t *testing.T) {
	var badPlugin BadPlugin

	conf := uconfig.New[f.Config](badPlugin)
	_, err := conf.Parse()

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "unsupported plugins. expecting a walker or visitor" {
		t.Errorf("Expected unsupported plugin error, got: %v", err)
	}
}

type FailingPluginWalker struct {
	plugins.Plugin
}

func (fp FailingPluginWalker) Walk(any) error {
	return errors.New("failed to walk")
}

func TestFailingPlugWalker(t *testing.T) {
	var failingPluginWalker FailingPluginWalker

	conf := uconfig.New[f.Config](failingPluginWalker)
	_, err := conf.Parse()

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "failed to walk" {
		t.Errorf("Expected failed to walk, got: %v", err)
	}
}

type FailingPluginVisitor struct {
	plugins.Plugin
}

func (fp FailingPluginVisitor) Visit(flat.Fields) error {
	return errors.New("failed to visit")
}

func TestFailingPlugVisitor(t *testing.T) {
	var failingPluginVisitor FailingPluginVisitor

	conf := uconfig.New[f.Config](failingPluginVisitor)
	_, err := conf.Parse()

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "failed to visit" {
		t.Errorf("Expected failed to visit, got: %v", err)
	}
}

// --- Extension tests ---

type testExtension struct {
	plugins []plugins.Plugin
	err     error
}

func (e *testExtension) Extend(ps []plugins.Plugin) error {
	e.plugins = ps
	return e.err
}
func (e *testExtension) Parse() error { return nil }

func TestExtensionReceivesPlugins(t *testing.T) {
	ext := &testExtension{}
	conf := uconfig.New[f.Config](ext)

	_, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(ext.plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(ext.plugins))
	}
	if ext.plugins[0] != ext {
		t.Fatal("expected extension to see itself in the plugin list")
	}
}

func TestExtensionError(t *testing.T) {
	ext := &testExtension{err: errors.New("extend failed")}
	conf := uconfig.New[f.Config](ext)

	_, err := conf.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "extend failed" {
		t.Fatalf("expected 'extend failed', got: %v", err)
	}
}

// --- Watch / Updater tests ---

type simpleConfig struct {
	Name string `default:"hello"`
}

// testUpdater is a fake Extension+Updater for testing Watch.
type testUpdater struct {
	ch chan struct{}
}

func (u *testUpdater) Extend([]plugins.Plugin) error { return nil }
func (u *testUpdater) Parse() error                  { return nil }

func (u *testUpdater) Updated(ctx context.Context) bool {
	select {
	case <-u.ch:
		return true
	case <-ctx.Done():
		return false
	}
}

func TestWatchNoUpdaters(t *testing.T) {
	conf := uconfig.New[simpleConfig]()

	var called bool
	err := conf.Watch(context.Background(), func(ctx context.Context, c *simpleConfig) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !called {
		t.Fatal("callback should be called once")
	}
}

func TestWatchNoUpdatersReturnsError(t *testing.T) {
	conf := uconfig.New[simpleConfig]()
	sentinel := errors.New("boom")

	err := conf.Watch(context.Background(), func(ctx context.Context, c *simpleConfig) error {
		return sentinel
	})

	if err != sentinel {
		t.Fatalf("expected sentinel, got %v", err)
	}
}

func TestWatchFnReturnsNilExits(t *testing.T) {
	updater := &testUpdater{ch: make(chan struct{}, 1)}
	conf := uconfig.New[simpleConfig](updater)

	err := conf.Watch(context.Background(), func(ctx context.Context, c *simpleConfig) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWatchCallbackError(t *testing.T) {
	updater := &testUpdater{ch: make(chan struct{}, 1)}
	conf := uconfig.New[simpleConfig](updater)

	sentinel := errors.New("callback failed")
	err := conf.Watch(context.Background(), func(ctx context.Context, c *simpleConfig) error {
		return sentinel
	})

	if err != sentinel {
		t.Fatalf("expected sentinel, got %v", err)
	}
}

func TestWatchContextCancel(t *testing.T) {
	updater := &testUpdater{ch: make(chan struct{})}
	conf := uconfig.New[simpleConfig](updater)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- conf.Watch(ctx, func(ctx context.Context, c *simpleConfig) error {
			<-ctx.Done()
			return nil
		})
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch didn't stop after cancel")
	}
}

func TestWatchReloadsOnUpdate(t *testing.T) {
	updater := &testUpdater{ch: make(chan struct{}, 1)}
	conf := uconfig.New[simpleConfig](updater)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls atomic.Int32
	started := make(chan struct{}, 5)

	go conf.Watch(ctx, func(ctx context.Context, c *simpleConfig) error {
		calls.Add(1)
		started <- struct{}{}
		<-ctx.Done()
		return nil
	})

	<-started
	if calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", calls.Load())
	}

	updater.ch <- struct{}{}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for reload")
	}

	if calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", calls.Load())
	}
}

func TestWatchMultipleUpdaters(t *testing.T) {
	u1 := &testUpdater{ch: make(chan struct{}, 1)}
	u2 := &testUpdater{ch: make(chan struct{}, 1)}
	conf := uconfig.New[simpleConfig](u1, u2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls atomic.Int32
	started := make(chan struct{}, 5)

	go conf.Watch(ctx, func(ctx context.Context, c *simpleConfig) error {
		calls.Add(1)
		started <- struct{}{}
		<-ctx.Done()
		return nil
	})

	<-started // initial

	// Signal from second updater only.
	u2.ch <- struct{}{}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for reload from second updater")
	}

	if calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", calls.Load())
	}
}

func TestWatchFnBlocksUntilChange(t *testing.T) {
	updater := &testUpdater{ch: make(chan struct{}, 1)}
	conf := uconfig.New[simpleConfig](updater)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls atomic.Int32
	started := make(chan struct{}, 5)
	done := make(chan error, 1)

	go func() {
		done <- conf.Watch(ctx, func(ctx context.Context, c *simpleConfig) error {
			n := calls.Add(1)
			started <- struct{}{}
			if n >= 2 {
				return nil // exit on second call
			}
			<-ctx.Done() // block until change
			return nil
		})
	}()

	<-started // initial call, blocking on ctx.Done

	// Signal a change.
	updater.ch <- struct{}{}

	select {
	case <-started: // second call
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for reload")
	}

	// Second call returns nil — Watch should exit.
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch didn't exit after fn returned nil")
	}
}
