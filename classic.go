package uconfig

import (
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/flag"
)

// Classic creates a uconfig manager with defaults,environment variables, and flags (in that order) and parses them right away.
func Classic(conf interface{}) (Config, error) {
	c, err := New(conf,
		defaults.New(),
		env.New(),
		flag.Standard(),
	)

	if err != nil {
		return c, err
	}
	return c, c.Parse()
}
