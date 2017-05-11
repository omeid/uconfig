// Package env provides flags integration for uconfig
package env

import (
	"os"

	"github.com/danverbraganza/varcaser/varcaser"
	"github.com/omeid/uconfig/flat"
)

const tag = "env"

// Envs is an env variable plugin.
type Envs interface {
	Visit(flat.Fields) error

	Parse() error
}

// New returns an EnvSet.
func New() Envs {
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

	return v.fields.Visit(func(f flat.Field) error {

		name, ok := f.Tag(tag)
		if name == "-" {
			return nil
		}

		if !ok || name == "" {
			name = f.Name()
			name = v.vc.String(name)
		}

		f.Meta()[tag] = "$" + name

		return nil
	})
}

func (v *visitor) Parse() error {

	return v.fields.Visit(func(f flat.Field) error {

		// Next block could use field.Meta and grab the tag name.
		/* start */
		name, ok := f.Tag(tag)
		if name == "-" {
			return nil
		}

		if !ok || name == "" {
			name = f.Name()
			name = v.vc.String(name)
		}

		/* end */
		value := os.Getenv(name)

		if value == "" {
			return nil
		}
		return f.Set(value)
	})
}
