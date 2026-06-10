package file

import (
	"os"
	"path/filepath"

	"github.com/omeid/uconfig/paths"
)

type absolute string

func (a absolute) Name() string     { return string(a) }
func (a absolute) Kind() paths.Kind { return paths.Absolute }
func (a absolute) Resolve() string  { return string(a) }

// Absolute returns a paths.One for a fixed absolute path.
func Absolute(path string) paths.One {
	return absolute(path)
}

type relative string

func (r relative) Name() string     { return string(r) }
func (r relative) Kind() paths.Kind { return paths.Relative }
func (r relative) Resolve() string {
	abs, err := filepath.Abs(string(r))
	if err != nil {
		return string(r)
	}
	return abs
}

// Relative returns a paths.One that resolves a relative path against
// the working directory at the time of the call.
func Relative(path string) paths.One {
	return relative(path)
}

type workspace string

func (w workspace) Name() string     { return string(w) }
func (w workspace) Kind() paths.Kind { return paths.Workspace }
func (w workspace) Resolve() string {
	dir, err := filepath.Abs(".")
	if err != nil {
		return ""
	}

	for {
		candidate := filepath.Join(dir, string(w))
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return ""
		}
		dir = parent
	}
}

// Workspace returns a paths.One that walks up the directory tree looking
// for a file at the given relative path (e.g. ".myapp/config").
// Returns the absolute path of the first match, or empty string if
// not found.
//
// The search always starts from the current working directory.
//
// This implements the common ancestor-directory search pattern used by
// tools like git (.git), eslint (.eslintrc), and similar.
func Workspace(name string) paths.One {
	return workspace(name)
}
