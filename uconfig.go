// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import (
	"fmt"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

var ErrUsage = plugins.ErrUsage

// Config is the config manager.
type Config interface {
	// Parse will call the parse method of all the added pluginss in the order
	// that the pluginss were registered, it will return early as soon as any
	// plugins fails.
	// You must call this before using the config value.
	Parse() error

	// Usage provides a simple usage message based on the meta data registered
	// by the pluginss.
	Usage()
}

// New returns a new Config. The conf must be a pointer to a struct.
func New(conf interface{}, ps ...plugins.Plugin) (Config, error) {
	fields, err := flat.View(conf)

	c := &config{
		conf:    conf,
		fields:  fields,
		plugins: make([]plugins.Plugin, 0, len(ps)),
	}

	if err != nil {
		return c, err
	}

	for _, plug := range ps {

		err := c.addPlugin(plug)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

type config struct {
	plugins []plugins.Plugin
	conf    interface{}
	fields  flat.Fields
}

func (c *config) addPlugin(plug plugins.Plugin) error {
	switch plug := plug.(type) {

	case plugins.Visitor:
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
