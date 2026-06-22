// Package file provides config file support for uconfig
package file

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/omeid/uconfig/plugins"
)

// Kind represents the type of a path resolver.
type Kind int

const (
	KindUnknown Kind = iota
	KindAbsolute
	KindRelative
	KindWorkspace
	KindStdin
)

func (k Kind) String() string {
	switch k {
	case KindAbsolute:
		return "absolute"
	case KindRelative:
		return "relative"
	case KindWorkspace:
		return "workspace"
	case KindStdin:
		return "stdin"
	default:
		return "unknown"
	}
}

// Source represents an abstraction for a configuration file source.
type Source interface {
	// Name returns the display name of the path (e.g., as written by the user).
	Name() string
	// Kind returns the kind of the path (e.g., Absolute).
	Kind() Kind
	// Resolve returns the actual absolute filesystem path to the file.
	Resolve() string
	// Open opens the source for reading.
	Open() (io.ReadCloser, error)
}

// Unmarshal is any function that maps the source bytes to the provided
// config.
type Unmarshal func(src []byte, v any) error

// Files represents a list of files to load.
type Files []struct {
	Path      Source
	Unmarshal Unmarshal
	Optional  bool
}

// Plugins constructs a slice of Plugin from the Files list.
func (f Files) Plugins() []plugins.Plugin {
	ps := make([]plugins.Plugin, 0, len(f))
	for _, file := range f {
		ps = append(ps, New(file.Path, file.Unmarshal, Optional(file.Optional)))
	}

	return ps
}



// Option configures a file walker.
type Option func(*walker)

// Optional configures whether the file plugin should fail if the source does not exist.
func Optional(b bool) Option {
	return func(w *walker) {
		w.optional = b
	}
}

// New returns a file plugin from a source and unmarshal function.
func New(source Source, unmarshal Unmarshal, opts ...Option) plugins.Plugin {
	w := &walker{
		Path:      source,
		Unmarshal: unmarshal,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

type walker struct {
	Path      Source
	Unmarshal Unmarshal
	optional  bool
	conf      any

	cachedReader io.ReadCloser // holds an already opened reader to avoid re-opening (e.g. for io.Reader sources)
}

// SourcePaths implements plugins.SourcePaths.
func (w *walker) SourcePaths() []plugins.SourcePath {
	if w.Path == nil {
		return nil
	}
	path := w.Path.Resolve()
	if path != "" && w.Path.Kind() != KindStdin && w.Path.Kind() != KindUnknown {
		return []plugins.SourcePath{{Name: w.Path.Name(), Path: path}}
	}
	return nil
}

// Usage implements plugins.Usage.
func (w *walker) Usage(fieldname string) (string, string) {
	if fieldname == "." && w.Path != nil {
		name := w.Path.Name()
		kind := w.Path.Kind()

		if kind != KindUnknown {
			return "Files", fmt.Sprintf("%-10s %s", kind.String()+":", name)
		}
		return "Files", name
	}
	return "", ""
}

func (w *walker) Walk(conf any) error {
	w.conf = conf

	if w.Path == nil {
		return nil
	}

	// For stdin, do the terminal check eagerly during Walk.
	if w.Path.Kind() == KindStdin {
		stat, err := os.Stdin.Stat()
		if err == nil {
			isTerminal := (stat.Mode() & os.ModeCharDevice) != 0
			isEmptyRegular := stat.Mode().IsRegular() && stat.Size() <= 0

			if isTerminal || isEmptyRegular {
				if w.optional {
					w.Path = nil // skip reading
				} else if isTerminal {
					return errors.New("stdin is required but no data was provided")
				}
			}
		}
		return nil
	}

	// Eagerly check file existence if not optional.
	// We don't want to fail if the source isn't an os.File.
	filepath := w.Path.Resolve()
	if filepath != "" && filepath != "-" && w.Path.Kind() != KindUnknown {
		_, err := os.Stat(filepath)
		if err != nil {
			if w.optional && os.IsNotExist(err) {
				return nil
			}
			return err
		}
	}

	return nil
}

func (w *walker) Parse() error {
	if w.Path == nil {
		return nil
	}

	var rc io.ReadCloser
	if w.cachedReader != nil {
		rc = w.cachedReader
		w.cachedReader = nil // Consume it
	} else {
		var err error
		rc, err = w.Path.Open()
		if err != nil {
			if w.optional && os.IsNotExist(err) {
				return nil
			}
			return err
		}
	}

	data, err := io.ReadAll(rc)
	_ = rc.Close()

	if err != nil {
		return err
	}

	err = w.Unmarshal(data, w.conf)
	if err != nil {
		filepath := w.Path.Resolve()
		filePathErr := errors.New(filepath)
		return errors.Join(filePathErr, err)
	}

	return nil
}
