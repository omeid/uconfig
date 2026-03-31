// Package file provides config file support for uconfig
package file

import (
	"errors"
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
		name:      filepath,
		filepath:  filepath,
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

// FilePaths returns the resolved filesystem paths from a list of
// plugins, filtering out non-file plugins. Paths are available
// after Walk has been called.
func FilePaths(ps []plugins.Plugin) []string {
	var paths []string
	for _, p := range ps {
		if w, ok := p.(*walker); ok && w.filepath != "" {
			paths = append(paths, w.filepath)
		}
	}
	return paths
}

// FileNames returns the display names of file paths from a list of
// plugins, filtering out non-file plugins. These are the names as
// provided by the user, not resolved absolute paths.
func FileNames(ps []plugins.Plugin) []string {
	var names []string
	for _, p := range ps {
		if w, ok := p.(*walker); ok && w.name != "" {
			names = append(names, w.name)
		}
	}
	return names
}

type walker struct {
	name      string        // display name (as the user wrote it)
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

	if w.src != nil {
		// Created via NewReader -- use the provided reader (one-shot).
		src = w.src
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

	if closer, ok := src.(io.Closer); ok {
		if err := closer.Close(); err != nil {
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
