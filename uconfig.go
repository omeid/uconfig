// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import (
	"errors"
	"fmt"
	"os"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

var ErrUsage = plugins.ErrUsage

// Config is the config manager.
type Config[C any] interface {
	// Parse will call the parse method of all the added pluginss in the order
	// that the pluginss were registered, it will return early as soon as any
	// plugins fails.
	// You must call this before using the config value.
	Parse() (*C, error)

	// Run calls parse and checks the error to see if usage was request
	// otherwise prints the error and usage and exists with os.Exit(1)
	Run() *C

	// Usage provides a simple usage message based on the meta data registered
	// by the pluginss.
	Usage()
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
		ret := 0
		if !usageRequest {
			fmt.Println(err)
		}
		c.Usage()
		os.Exit(ret)
	}

	return conf
}
