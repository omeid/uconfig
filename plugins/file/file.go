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
	Path      string
	Unmarshal Unmarshal
	Optional  bool
}

// Plugins constructs a slice of Plugin from the Files list of
// paths and unmarshal functions.
func (f Files) Plugins() []plugins.Plugin {
	ps := make([]plugins.Plugin, 0, len(f))
	for _, f := range f {
		fp := New(
			f.Path,
			f.Unmarshal,
			Config{Optional: f.Optional},
		)

		ps = append(ps, fp)
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
		filepath:  path,
		unmarshal: unmarshal,
		optional:  config.Optional,
	}

	return plug
}

// FilePaths returns the file paths from a list of plugins,
// filtering out non-file plugins.
func FilePaths(ps []plugins.Plugin) []string {
	var paths []string
	for _, p := range ps {
		if w, ok := p.(*walker); ok && w.filepath != "" {
			paths = append(paths, w.filepath)
		}
	}
	return paths
}

type walker struct {
	filepath  string
	src       io.Reader // only set when created via NewReader
	conf      any
	unmarshal Unmarshal
	optional  bool

	err error
}

func (w *walker) Walk(conf any) error {
	w.conf = conf

	// Check file exists early (for non-optional files).
	if w.src == nil {
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

var ErrEncodingFailed = errors.New("failed to decode file")

func (w *walker) Parse() error {
	if w.err != nil {
		return w.err
	}

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
		defer f.Close()
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
