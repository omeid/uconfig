package file

import (
	"os"
	"path/filepath"
)

// Path pairs a display name with a lazy path resolver.
// The Name is shown in usage output; Resolve is called during
// Walk to obtain the actual filesystem path.
type Path struct {
	Name    string
	Resolve func() string
}

// Absolute returns a Path for a fixed absolute path.
func Absolute(path string) Path {
	return Path{Name: "absolute:  " + path, Resolve: func() string { return path }}
}

// Relative returns a Path that resolves a relative path against
// the working directory at the time of the call.
func Relative(path string) Path {
	return Path{
		Name: "relative:  " + path,
		Resolve: func() string {
			abs, err := filepath.Abs(path)
			if err != nil {
				return path
			}
			return abs
		},
	}
}

// Workspace returns a Path that walks up the directory tree looking
// for a file at the given relative path (e.g. ".myapp/config").
// Returns the absolute path of the first match, or empty string if
// not found.
//
// The search always starts from the current working directory.
//
// This implements the common ancestor-directory search pattern used by
// tools like git (.git), eslint (.eslintrc), and similar.
func Workspace(name string) Path {
	return Path{
		Name: "workspace: " + name,
		Resolve: func() string {
			dir, err := filepath.Abs(".")
			if err != nil {
				return ""
			}

			for {
				candidate := filepath.Join(dir, name)
				if _, err := os.Stat(candidate); err == nil {
					return candidate
				}
				parent := filepath.Dir(dir)
				if parent == dir {
					return "" // reached filesystem root
				}
				dir = parent
			}
		},
	}
}
