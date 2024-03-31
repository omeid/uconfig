package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/omeid/uconfig/plugins"
)

type UnmarshalOptions map[string]Unmarshal

var ErrFileExtNotSupported = errors.New("File extension not supported.")

// NewMutli returns a multi unmarshal plugin that can decode the file from path
// using various Unmarshal functions provided in unmarshal map.
// This is usually used as a second stage to load configurations based on a flag
// or configuration value.
func NewMulti(path string, unmarshalOptions UnmarshalOptions, optional bool) plugins.Plugin {

	plug := &multiWalker{
		filepath:         path,
		unmarshalOptions: unmarshalOptions,
	}

	ext := filepath.Ext(path)

	_, ok := unmarshalOptions[ext]
	if !ok {
		plug.err = ErrFileExtNotSupported
		return plug
	}

	src, err := os.Open(path)

	if optional && os.IsNotExist(err) {
		return plug
	}

	if err != nil {
		plug.err = err
		return plug
	}

	plug.src = src
	return plug
}

type multiWalker struct {
	filepath         string
	src              io.Reader
	conf             interface{}
	unmarshalOptions map[string]Unmarshal

	err error
}

func (v *multiWalker) Walk(conf interface{}) error {
	if v.err != nil {
		return v.err
	}

	v.conf = conf
	return v.err
}

func (v *multiWalker) Parse() error {

	if v.err != nil {
		return v.err
	}

	if v.src == nil {
		return nil
	}

	src, err := io.ReadAll(v.src)
	if err != nil {
		return err
	}

	if closer, ok := v.src.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	ext := filepath.Ext(v.filepath)
	unmarshal, ok := v.unmarshalOptions[ext]
	if !ok {
		return ErrFileExtNotSupported
	}

	return unmarshal(src, v.conf)
}
