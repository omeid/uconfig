package fileset

import (
	"os"
	"path/filepath"

	"github.com/omeid/uconfig/plugins/file"
)

// Set represents a collection of file sources matched by a pattern.
type Set interface {
	Name() string
	Kind() file.Kind
	Resolve() ([]file.Source, error)
}

type absolute string

func (a absolute) Name() string    { return string(a) }
func (a absolute) Kind() file.Kind { return file.KindAbsolute }
func (a absolute) Resolve() ([]file.Source, error) {
	matches, err := filepath.Glob(string(a))
	if err != nil {
		return nil, err
	}
	var sources []file.Source
	for _, match := range matches {
		sources = append(sources, file.Absolute(match))
	}
	return sources, nil
}

// Absolute returns a fileset sourced from an absolute glob pattern.
func Absolute(pattern string) Set {
	return absolute(pattern)
}

type relative string

func (r relative) Name() string    { return string(r) }
func (r relative) Kind() file.Kind { return file.KindRelative }
func (r relative) Resolve() ([]file.Source, error) {
	abs, err := filepath.Abs(string(r))
	if err != nil {
		return nil, err
	}
	matches, err := filepath.Glob(abs)
	if err != nil {
		return nil, err
	}
	var sources []file.Source
	for _, match := range matches {
		sources = append(sources, file.Absolute(match))
	}
	return sources, nil
}

// Relative returns a fileset sourced from a relative glob pattern.
func Relative(pattern string) Set {
	return relative(pattern)
}

type workspace string

func (w workspace) Name() string    { return string(w) }
func (w workspace) Kind() file.Kind { return file.KindWorkspace }
func (w workspace) Resolve() ([]file.Source, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		candidate := filepath.Join(dir, string(w))
		matches, err := filepath.Glob(candidate)
		if err == nil && len(matches) > 0 {
			var sources []file.Source
			for _, match := range matches {
				sources = append(sources, file.Absolute(match))
			}
			return sources, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, nil
}

// Workspace returns a fileset sourced from a glob pattern searched upwards in the workspace.
func Workspace(pattern string) Set {
	return workspace(pattern)
}
