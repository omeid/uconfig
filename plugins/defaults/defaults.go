// Package defaults provides flags integration for uconfig
package defaults

import (
	"github.com/danverbraganza/varcaser/varcaser"
	"github.com/omeid/uconfig/flat"
)

const tag = "default"

// Defaults is an env variable plugin.
type Defaults interface {
	Visit(flat.Fields) error

	Parse() error
}

// New returns an EnvSet.
func New() Defaults {
	return &visitor{
		vc: varcaser.Caser{
			From: varcaser.UpperCamelCase,
			To:   varcaser.ScreamingSnakeCase,
		},
	}
}

type visitor struct {
	vc     varcaser.Caser
	fields flat.Fields
}

func (v *visitor) Visit(f flat.Fields) error {

	v.fields = f

	return nil
}

func (v *visitor) Parse() error {

	return v.fields.Visit(func(f flat.Field) error {

		value, ok := f.Tag(tag)
		if !ok {
			return nil
		}

		f.Meta()[tag] = value
		return f.Set(value)
	})
}
