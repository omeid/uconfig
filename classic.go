package uconfig

import (
	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/omeid/uconfig/plugins/flag"
)

// Files represents a set of file paths and the appropriate unmarshaller function.
type Files = file.Files

// Classic creates a uconfig manager with defaults,environment variables,
// and flags (in that order) and optionally file loaders based on the provided
// PluginProvider (e.g. file.Files or watchfile.Files) and parses them right away.
func Classic[C any](files PluginProvider, userPlugins ...plugins.Plugin) Config[C] {
	// almost a duplicate of Load, but due to the order of things, not worth abstracting for a few lines.
	ps := make([]plugins.Plugin, 0, 3+len(userPlugins))
	// first defaults
	ps = append(ps, defaults.New())
	// then files
	if files != nil {
		ps = append(ps, files.Plugins()...)
	}
	// then any user plugins, often just _secret_.
	ps = append(ps, userPlugins...)

	// followed by envs
	ps = append(ps, env.New())

	// and lastly flags.
	ps = append(ps, flag.Standard())

	return New[C](ps...)
}
