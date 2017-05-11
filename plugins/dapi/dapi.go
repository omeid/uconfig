package dapi

import (
	"errors"
	"strings"

	"github.com/danverbraganza/varcaser/varcaser"
	"github.com/omeid/uconfig/flat"
)

const tag = "dapi"

// DAPI plugin for uconfig
type DAPI interface {
	Visit(flat.Fields) error
	Parse() error
}

// New returns an EnvSet.
func New(base string) DAPI {
	return &visitor{
		base: base,
		vc: varcaser.Caser{
			From: varcaser.UpperCamelCase,
			To:   varcaser.LowerCamelCase,
		},
	}
}

type visitor struct {
	vc     varcaser.Caser
	fields flat.Fields

	base string
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

		return checkFile(file)

	})
	return nil
}

func (v *visitor) Parse() error {

	return v.fields.Visit(func(f flat.Field) error {
		tag, ok := f.Tag(tag)

		if !ok {
			return nil
		}

		file, field, err := splitTag(tag)

		value, err := readField(v.base, file, field)

		if err != nil {
			return err
		}

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

func checkFile(path string) error {
	return nil
}

func readField(base string, file string, field string) (string, error) {
	return "", nil
}
