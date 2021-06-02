// Package file provides config file support for uconfig
package file

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/omeid/uconfig/plugins"
)

// Files represents a set of file paths and the appropriate
// unmarshal function for the given file.
type Files map[string]Unmarshal

// Plugins constructs a slice of Plugin from the Files list of
// paths and unmarshal functions.
func (f Files) Plugins(config Config) []plugins.Plugin {
	ps := make([]plugins.Plugin, 0, len(f))
	for path, unmarshal := range f {
		ps = append(ps, New(path, unmarshal, config))
	}

	return ps
}

// Unmarshal is any function that maps the source bytes to the provided
// config.
type Unmarshal func(src []byte, v interface{}) error

// NewReader returns a uconfig plugin that unmarshals the content of
// the provided io.Reader into the config using the provided unmarshal
// function. The src will be closed if it is an io.Closer.
func NewReader(src io.Reader, unmarshal Unmarshal) plugins.Plugin {
	return &walker{
		src:       src,
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

	src, err := os.Open(path)

	if err != nil && config.Optional && os.IsNotExist(err) {
		return &walker{
			unmarshal: unmarshal,
		}
	}

	return &walker{
		src:       src,
		unmarshal: unmarshal,
		err:       err,
	}
}

type walker struct {
	src       io.Reader
	conf      interface{}
	unmarshal Unmarshal

	err error
}

func (v *walker) Walk(conf interface{}) error {
	if v.err != nil {
		return v.err
	}

	v.conf = conf
	return v.err
}

func (v *walker) Parse() error {

	if v.err != nil {
		return v.err
	}

	if v.src == nil {
		return nil
	}

	src, err := ioutil.ReadAll(v.src)
	if err != nil {
		return err
	}

	if closer, ok := v.src.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	return v.unmarshal(src, v.conf)
}
