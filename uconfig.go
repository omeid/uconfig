// Package uconfig provides advanced command line flags supporting defaults, env vars, and config structs.
package uconfig

import "github.com/omeid/uconfig/flat"

// Plugin is the common interface for all plugins.
type Plugin interface {
	Parse() error
}

// Walker is the interface for plugins that take the whole config, like file loaders.
type Walker interface {
	Plugin

	Walk(interface{}) error
}

// WalkerFunc is a helper type that turns a Walk function into a Walker.
// type WalkerFunc func(interface{}) error

// Walk implements Walker for WalkerFunc
// func (wf WalkerFunc) Walk(conf interface{}) error { return wf(conf) }

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
	Visitor(Visitor) error
	Walker(Walker) error

	// Must be called after Visitor and Walkers are added.
	Usage()
	Parse() error
}

// New returns a new Config.
func New(conf interface{}) (Config, error) {
	fields, err := flat.View(conf)
	if err != nil {
		return nil, err
	}

	c := &config{
		conf:   conf,
		fields: fields,
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
