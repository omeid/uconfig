// Package file provides config file support for uconfig
package file

import (
	"io"
	"io/ioutil"
	"os"
)

// Unmarshal is any function that maps the source bytes to the provided
// config.
type Unmarshal func(src []byte, v interface{}) error

// File is a file loader plugin for uconfig
type File interface {
	Walk(conf interface{}) error
	Parse() error
}

// NewReader returns a uconfig plugin that unmarshals the content of
// the provided io.Reader into the config using the provided unmarshal
// function. The src will be closed if it is an io.Closer.
func NewReader(src io.Reader, unmarshal Unmarshal) File {
	return &walker{
		src:       src,
		unmarshal: unmarshal,
	}

}

// New returns an EnvSet.
func New(path string, unmarshal Unmarshal) File {
	src, err := os.Open(path)
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
	v.conf = conf

	return v.err
}

func (v *walker) Parse() error {

	src, err := ioutil.ReadAll(v.src)
	if err != nil {
		return err
	}

	if closer, ok := v.src.(io.Closer); ok {
		closer.Close()
	}

	return v.unmarshal(src, v.conf)
}
