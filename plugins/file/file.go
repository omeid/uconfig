// Package file provides config file support for uconfig
package file

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/omeid/uconfig/plugins"
)

// Files represents a set of file paths and the appropriate
// unmarshal function for the given file.
type Files []struct {
	Path      Path
	Unmarshal Unmarshal
	Optional  bool
}

// Plugins constructs a slice of Plugin from the Files list of
// paths and unmarshal functions.
func (f Files) Plugins() []plugins.Plugin {
	ps := make([]plugins.Plugin, 0, len(f))
	for _, f := range f {
		ps = append(ps, &walker{
			name:      f.Path.Name,
			kind:      f.Path.Kind,
			resolve:   f.Path.Resolve,
			unmarshal: f.Unmarshal,
			optional:  f.Optional,
		})
	}

	return ps
}

// Unmarshal is any function that maps the source bytes to the provided
// config.
type Unmarshal func(src []byte, v any) error

// NewReader returns a uconfig plugin that unmarshals the content of
// the provided io.Reader into the config using the provided unmarshal
// function. The src will be closed if it is an io.Closer.
func NewReader(src io.Reader, filepath string, unmarshal Unmarshal) plugins.Plugin {
	return &walker{
		src:       src,
		filepath:  filepath,
		name:      filepath,
		unmarshal: unmarshal,
	}
}

// Config describes the options required for a file.
type Config struct {
	// indicates if a file that does not exist should be ignored.
	Optional bool
}

// New returns a file plugin.
func New(path string, unmarshal Unmarshal, config Config) plugins.Plugin {
	plug := &walker{
		name:      path,
		filepath:  path,
		unmarshal: unmarshal,
		optional:  config.Optional,
	}

	return plug
}

// SourcePaths implements plugins.SourcePaths.
func (w *walker) SourcePaths() []plugins.SourcePath {
	if w.filepath != "" {
		return []plugins.SourcePath{{Name: w.name, Path: w.filepath}}
	}
	return nil
}

// Usage implements plugins.Usage.
func (w *walker) Usage(fieldname string) (string, string) {
	if fieldname == "." {
		if w.kind != Unknown {
			return "Files", fmt.Sprintf("%-10s %s", w.kind.String()+":", w.name)
		}
		return "Files", w.name
	}
	return "", ""
}

type walker struct {
	name      string        // display name (as the user wrote it)
	kind      Kind          // path kind (absolute, relative, workspace)
	filepath  string        // resolved absolute path (set during Walk)
	resolve   func() string // lazy resolver (from Path.Resolve)
	src       io.Reader     // only set when created via NewReader
	conf      any
	unmarshal Unmarshal
	optional  bool
}

func (w *walker) Walk(conf any) error {
	w.conf = conf

	// Lazy path resolution (e.g. Workspace, Relative).
	if w.resolve != nil && w.filepath == "" {
		w.filepath = w.resolve()
	}

	// Check file exists early (for non-optional files).
	if w.src == nil && w.filepath != "" {
		_, err := os.Stat(w.filepath)
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
	var src io.Reader
	var closeSrc io.Closer

	if w.src != nil {
		// Created via NewReader -- use the provided reader (one-shot).
		src = w.src
		if c, ok := src.(io.Closer); ok {
			closeSrc = c
		}
		w.src = nil // consumed
	} else {
		// Created via New -- open the file fresh each time.
		f, err := os.Open(w.filepath)
		if err != nil {
			if w.optional && os.IsNotExist(err) {
				return nil
			}
			return err
		}
		defer f.Close() //nolint:errcheck // read-only
		src = f
	}

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	if closeSrc != nil {
		if err := closeSrc.Close(); err != nil {
			return err
		}
	}

	err = w.unmarshal(data, w.conf)
	if err != nil {
		filePath := errors.New(w.filepath)
		return errors.Join(filePath, err)
	}

	return nil
}
