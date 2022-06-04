// Package plugins describes the uconfig provider interface.
// it exists to enable uconfig.Classic without circular deps.
package plugins

import (
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

	Walk(interface{}) error
}

// Visitor is the interface for providers that require a flat view
// of the config, like flags, env vars
type Visitor interface {
	Plugin

	Visit(flat.Fields) error
}

var tags = map[string]string{}

// ErrUsage is returned when user has request usage message
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
