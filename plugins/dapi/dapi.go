// Package dapi provides Kubernetes DownwardAPI ini config support for uconfig.
package dapi

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	"github.com/omeid/uconfig/flat"
)

const tag = "dapi"

// DAPI plugin for uconfig
type DAPI interface {
	Visit(flat.Fields) error
	Parse() error
}

// New returns DAPI plugin for uConfig that will load ini files from the provided base location. Please note that DAPI only works with explicitly tagged fields, the tags are in the form `dapi:"file_name:attribute"` where file_name is a file under base and attribute is the key expected in base/file_name.
func New(base string) DAPI {
	return &visitor{
		base: base,
	}
}

type visitor struct {
	fields flat.Fields

	base string

	files []string
}

func (v *visitor) Visit(f flat.Fields) error {

	v.fields = f

	f.Visit(func(f flat.Field) error {
		tag, ok := f.Tag(tag)

		if !ok {
			return nil
		}

		file, _, err := splitTag(tag)

		if err != nil {
			return err
		}

		f.Meta()[tag] = tag

		v.files = append(v.files, file)
		return nil

	})
	return nil
}

func (v *visitor) Parse() error {

	files := map[string]*ini.File{}

	for _, f := range v.files {
		cfg, err := ini.InsensitiveLoad(filepath.Join(v.base, f))
		if err != nil {
			return err
		}

		files[f] = cfg
	}

	return v.fields.Visit(func(f flat.Field) error {
		tag, ok := f.Tag(tag)

		if !ok {
			return nil
		}

		file, field, err := splitTag(tag)

		if err != nil {
			return err
		}

		value := files[file].Section("").Key(field).String()
		return f.Set(value)
	})

}

func splitTag(tag string) (string, string, error) {
	segs := strings.Split(tag, ":")

	if len(segs) != 2 {
		return "", "", errors.New("invalid dapi tag. Expecting `file:field`")
	}

	return segs[0], segs[1], nil
}
