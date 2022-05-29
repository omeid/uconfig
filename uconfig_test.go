package uconfig_test

import (
	"errors"
	"testing"

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

	config := f.Config{}

	_, err := uconfig.New(&config, badPlugin)

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "Unsupported plugins. Expecting a Walker or Visitor" {
		t.Errorf("Expected unsupported plugin error, got: %v", err)
	}
}

type FailingPluginWalker struct {
	plugins.Plugin
}

func (fp FailingPluginWalker) Walk(interface{}) error {
	return errors.New("failed to walk")
}

func TestFailingPlugWalker(t *testing.T) {

	var failingPluginWalker FailingPluginWalker

	config := f.Config{}

	_, err := uconfig.New(&config, failingPluginWalker)

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

	config := f.Config{}

	_, err := uconfig.New(&config, failingPluginVisitor)

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "failed to visit" {
		t.Errorf("Expected failed to visit, got: %v", err)
	}
}
