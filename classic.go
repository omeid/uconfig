package uconfig

import (
	"os"

	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/flag"
)

// Classic creates a uconfig manager with defaults,environment variables, and flags (in that order) and parses them right away.
//
//
//
func Classic(conf interface{}) error {
	c, err := New(conf)
	if err != nil {
		return err
	}
	fs := flag.New(os.Args[0], flag.ContinueOnError, os.Args[1:])
	c.Visitor(defaults.New())
	c.Visitor(env.New())
	c.Visitor(fs)
	return c.Parse()
}
