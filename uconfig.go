// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import (
	"fmt"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

// Config is the config manager.
type Config interface {
	// AddPlugin adds a Visitor or Walker plugin to the config.
	AddPlugin(plugins.Plugin) error

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

	c := &config{
		conf:   conf,
		fields: fields,
	}

	if err != nil {
		return c, err
	}

	for _, plug := range ps {

		err := c.AddPlugin(plug)
		if err != nil {
			return nil, err
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

func (c *config) AddPlugin(plug plugins.Plugin) error {
	switch plug := plug.(type) {

	case plugins.Visitor:
		// disable the std flag usage
		if plug, ok := plug.(canSetUsage); ok {
			plug.SetUsage(func() {})
		}

		err := plug.Visit(c.fields)
		if err != nil {
			return err
		}

	case plugins.Walker:
		err := plug.Walk(c.conf)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unsupported plugins. Expecting a Walker or Visitor")
	}

	c.plugins = append(c.plugins, plug)
	return nil
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
