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

// New returns an EnvSet.
func New(path string, unmarshal Unmarshal, config Config) plugins.Plugin {
	plug := &walker{
		filepath:  path,
		unmarshal: unmarshal,
	}

	src, err := os.Open(path)

	if err == nil {
		plug.src = src
	}

	if config.Optional && os.IsNotExist(err) {
		err = nil
	}

	plug.err = err

	return plug
}

type Plugin interface {
	plugins.Plugin
	FilePath() string
}

type walker struct {
	filepath  string
	src       io.Reader
	conf      any
	unmarshal Unmarshal

	err error
}

func (w *walker) FilePath() string {
	return w.filepath
}

func (w *walker) Walk(conf any) error {
	if w.err != nil {
		return w.err
	}

	w.conf = conf
	return w.err
}

var ErrEncodingFailed = errors.New("failed to decode file")

func (w *walker) Parse() error {
	if w.err != nil {
		return w.err
	}

	if w.src == nil {
		return nil
	}

	src, err := io.ReadAll(w.src)
	if err != nil {
		return err
	}

	if closer, ok := w.src.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	err = w.unmarshal(src, w.conf)
	if err != nil {
		filePath := errors.New(w.filepath)
		return errors.Join(filePath, err)
	}

	return nil
}
