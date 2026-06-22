package uconfig

import (
	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
)

// Load creates a uconfig manager with defaults,environment variables,
// and optionally file loaders based on the provided PluginProvider.
func Load[C any](files PluginProvider, userPlugins ...plugins.Plugin) Config[C] {
	ps := make([]plugins.Plugin, 0, 2+len(userPlugins))

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

	return New[C](ps...)
}
