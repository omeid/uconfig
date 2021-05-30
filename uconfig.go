// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import (
	"fmt"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

// Config is the config manager.
type Config interface {
	// Visitor adds a visitor plugins, Config invokes the plugins Visit method
	// right away with a flat view of the underlying config struct.
	Visitor(plugins.Visitor) error
	// Walker adds a walker plugins, Config invokes the plugins Walk method
	// right away with the underlying config struct.
	Walker(plugins.Walker) error

	// Must be called after Visitor and Walkers are added.
	// Parse will call the parse method of all the added pluginss in the order
	// that the pluginss were registered, it will return early as soon as any
	// plugins fails.
	Parse() error

	// Usage provides a simple usage message based on the meta data registered
	// by the pluginss.
	Usage()
}

// New returns a new Config. The conf must be a pointer to a struct.
func New(conf interface{}, ps ...plugins.Plugin) (Config, error) {
	fields, err := flat.View(conf)
	if err != nil {
		return nil, err
	}

	c := &config{
		conf:   conf,
		fields: fields,
	}

	for _, plug := range ps {
		switch plug := plug.(type) {
		case plugins.Visitor:
			err := c.Visitor(plug)
			if err != nil {
				return c, err
			}

		case plugins.Walker:
			err := c.Walker(plug)
			if err != nil {
				return c, err
			}
		default:
			return nil, fmt.Errorf("Unsupported plugins. Expecting a Walker or Visitor")
		}
	}

	return c, nil
}

type config struct {
	plugins []plugins.Plugin
	conf    interface{}
	fields  flat.Fields
}

type canSetUsage interface {
	SetUsage(func())
}

func (c *config) Visitor(v plugins.Visitor) error {

	// disable the std flag usage
	if v, ok := v.(canSetUsage); ok {
		v.SetUsage(func() {})
	}

	err := v.Visit(c.fields)
	if err != nil {
		return err
	}
	c.addPlugin(v)

	return nil
}

func (c *config) Walker(w plugins.Walker) error {
	err := w.Walk(c.conf)
	if err != nil {
		return err
	}
	c.addPlugin(w)
	return nil
}

func (c *config) addPlugin(p plugins.Plugin) {
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
