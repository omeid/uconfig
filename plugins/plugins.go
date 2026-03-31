// Package plugins describes the uconfig provider interface.
// It exists to enable uconfig.Classic without circular deps.
package plugins

import (
	"context"
	"errors"
	"log"
	"runtime"

	"github.com/omeid/uconfig/flat"
)

// Plugin is the common interface for all uConfig providers.
type Plugin interface {
	Parse() error
}

// Walker is the interface for providers that take the whole
// config, like file loaders.
type Walker interface {
	Plugin

	Walk(any) error
}

// Visitor is the interface for providers that require a flat view
// of the config, like flags, env vars
type Visitor interface {
	Plugin

	Visit(flat.Fields) error
}

// Extension is the interface for plugins that need access to the full
// plugin list. Like all plugins, Extensions are set up in registration
// order — place them after any plugins they need to inspect (e.g. after
// file plugins so that paths are resolved).
type Extension interface {
	Plugin

	Extend([]Plugin) error
}

// Updater is an optional interface for plugins that can detect
// when their backing source has changed. It is used by Watch
// to trigger re-parsing. Any plugin type (Walker, Visitor, or
// Extension) can additionally implement Updater.
type Updater interface {
	// Updated blocks until the plugin's source has changed or ctx is done.
	// Returns true if the source changed, false otherwise.
	Updated(ctx context.Context) bool
}

var tags = map[string]string{}

// ErrUsage is returned when the user has requested a usage message
// via some plugin, mostly flags.
var ErrUsage = errors.New("uconfig: usage request")

// RegisterTag allows providers to ensure their tag is unique.
// they must call this function from an init.
func RegisterTag(name string) {
	if pkg, exists := tags[name]; exists {
		log.Panicf("tag %s already registered by %s", name, pkg)
	}

	pc, _, _, _ := runtime.Caller(1)
	tags[name] = runtime.FuncForPC(pc).Name()
}
