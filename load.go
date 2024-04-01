package uconfig

import (
	"errors"
	"os"

	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/file"
)

// UnmarshalOptions represents a set of file paths and the appropriate unmarshaller function.
type UnmarshalOptions = file.UnmarshalOptions

// Load creates a uconfig manager with defaults,environment variables,
// and optionally file loaders based on the provided
// Files map and parses them right away.
func Load(conf interface{}, files Files, userPlugins ...plugins.Plugin) (Config, error) {

	fps := files.Plugins()

	ps := make([]plugins.Plugin, 0, len(fps)+3+len(userPlugins))

	// first defaults
	ps = append(ps, defaults.New())
	// then files
	ps = append(ps, fps...)
	// then any user pugins, often just _secret_.
	ps = append(ps, userPlugins...)

	// followed by envs
	ps = append(ps, env.New())

	c, err := New(conf, ps...)

	if err != nil {
		return c, err
	}

	err = c.Parse()
	if errors.Is(err, ErrUsage) {
		c.Usage()
		os.Exit(0)
	}

	return c, err
}
