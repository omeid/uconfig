package uconfig

import (
	"errors"
	"os"

	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/omeid/uconfig/plugins/flag"
)

// Files represents a set of file paths and the appropriate unmarshaller function.
type Files = file.Files

func classicConfig(conf interface{}, files Files, userPlugins []plugins.Plugin, stripCommand bool) (*config, error) {
	fps := files.Plugins()

	ps := make([]plugins.Plugin, 0, len(fps)+3+len(userPlugins))

	// first defaults
	ps = append(ps, defaults.New())
	// then files
	ps = append(ps, fps...)
	// followed by env and flags
	ps = append(ps, env.New())

	stripArgs := 1
	if stripCommand && len(os.Args) > 1 {
		stripArgs = 2
	}

	ps = append(ps, flag.New(os.Args[0], flag.ContinueOnError, os.Args[stripArgs:]))

	// then any user pugins, often just _secret_.
	ps = append(ps, userPlugins...)

	return newConfig(conf, ps)
}

// Classic creates a uconfig manager with defaults,environment variables,
// and flags (in that order) and optionally file loaders based on the provided
// Files map and parses them right away.
func Classic(conf interface{}, files Files, userPlugins ...plugins.Plugin) (Config, error) {

	c, err := classicConfig(conf, files, userPlugins, false)

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
