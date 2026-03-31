// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

var ErrUsage = plugins.ErrUsage

// PluginProvider is implemented by types that can provide plugins.
// Both file.Files and watchfile.Files implement this interface,
// allowing Classic and Load to accept either.
type PluginProvider interface {
	Plugins() []plugins.Plugin
}

// Config is the config manager.
type Config[C any] interface {
	// Parse will call the parse method of all the added plugins in the order
	// they were registered. It returns early as soon as any plugin fails.
	// You must call this before using the config value.
	Parse() (*C, error)

	// Run calls Parse and checks the error to see if usage was requested,
	// otherwise prints the error and usage and exits with os.Exit(1).
	Run() *C

	// Usage provides a simple usage message based on the meta data registered
	// by the plugins.
	Usage()

	// Watch calls Parse for the initial configuration, then calls fn.
	// When any plugin that implements Updater signals a change, fn's
	// context is cancelled, the config is re-parsed, and fn is called
	// again with the new value.
	// If no plugins implement Updater, fn is called once.
	//
	// fn should block (e.g. <-ctx.Done()) to stay alive until a
	// config change. When fn returns, Watch exits with fn's error.
	Watch(ctx context.Context, fn func(ctx context.Context, c *C) error) error
}

// New returns a new Config. The conf must be a pointer to a struct.
func New[C any](ps ...plugins.Plugin) Config[C] {
	conf := new(C)
	fields, err := flat.View(conf)

	return &config[C]{
		err:     err,
		conf:    conf,
		fields:  fields,
		plugins: ps,
	}
}

type config[C any] struct {
	plugins []plugins.Plugin
	conf    *C
	fields  flat.Fields

	err error // lazy error
}

func (c *config[C]) Parse() (*C, error) {
	if c.err != nil {
		return nil, c.err
	}

	// first setup plugins.
	for _, plug := range c.plugins {
		switch plug := plug.(type) {

		case plugins.Visitor:
			err := plug.Visit(c.fields)
			if err != nil {
				return nil, err
			}

		case plugins.Walker:
			err := plug.Walk(c.conf)
			if err != nil {
				return nil, err
			}

		case plugins.Extension:
			err := plug.Extend(c.plugins)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unsupported plugins. expecting a walker or visitor")
		}
	}

	for _, p := range c.plugins {

		err := p.Parse()
		if err != nil {
			return nil, err
		}
	}

	return c.conf, nil
}

func (c *config[C]) Run() *C {
	conf, err := c.Parse()
	if err != nil {
		usageRequest := errors.Is(err, ErrUsage)
		ret := 1
		if usageRequest {
			ret = 0
		} else {
			fmt.Println(err)
		}
		c.Usage()
		os.Exit(ret)
	}

	return conf
}
