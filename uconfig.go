// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import (
	"fmt"

	"github.com/omeid/uconfig/flat"
)

// Plugin is the common interface for all plugins.
type Plugin interface {
	Parse() error
}

// Walker is the interface for plugins that take the whole config, like file loaders.
type Walker interface {
	Plugin

	Walk(interface{}) error
}

// Visitor is the interface for plugins that require a flat view of the config, like flags, env vars
type Visitor interface {
	Plugin

	Visit(flat.Fields) error
}

// VisitorFunc is a helper type that turns a Visitor function into a Visitor.
// type VisitorFunc func(flat.Field) error

// Visit implements Visitor for VisitorFunc
// func (vf VisitorFunc) Visit(f flat.Field) error { return vf(f) }

// Config is the config manager.
type Config interface {
	// Visitor adds a visitor plugin, Config invokes the plugins Visit method
	// right away with a flat view of the underlying config struct.
	Visitor(Visitor) error
	// Walker adds a walker plugin, Config invokes the plugins Walk method
	// right away with the underlying config struct.
	Walker(Walker) error

	// Must be called after Visitor and Walkers are added.
	// Parse will call the parse method of all the added plugins in the order
	// that the plugins were registered, it will return early as soon as any
	// plugin fails stops calling parse on plugins.
	Parse() error

	// Usage provides a simple usage message based on the meta data registered
	// by the plugins.
	Usage()
}

// New returns a new Config. The conf must be a pointer to a struct.
func New(conf interface{}, plugins ...Plugin) (Config, error) {
	fields, err := flat.View(conf)
	if err != nil {
		return nil, err
	}

	c := &config{
		conf:   conf,
		fields: fields,
	}

	for _, plugin := range plugins {
		switch plugin := plugin.(type) {
		case Visitor:
			err := c.Visitor(plugin)
			if err != nil {
				return c, err
			}

		case Walker:
			err := c.Walker(plugin)
			if err != nil {
				return c, err
			}
		default:
			return nil, fmt.Errorf("Unsupported Plugin. Expecting a Walker or Visitor")
		}
	}

	return c, nil
}

type config struct {
	plugins []Plugin
	conf    interface{}
	fields  flat.Fields
}

type canSetUsage interface {
	SetUsage(func())
}

func (c *config) Visitor(v Visitor) error {

	// A special case for standard library flag plugin.
	if v, ok := v.(canSetUsage); ok {
		v.SetUsage(c.Usage)
	}

	err := v.Visit(c.fields)
	if err != nil {
		return err
	}
	c.plug(v)

	return nil
}
func (c *config) Walker(w Walker) error {
	err := w.Walk(c.conf)
	if err != nil {
		return err
	}
	c.plug(w)
	return nil
}

func (c *config) plug(p Plugin) {
	c.plugins = append(c.plugins, p)
}

func (c *config) Parse() error {
	for _, p := range c.plugins {

		err := p.Parse()
		if err != nil {
			return err
		}
	}

	return nil
}
