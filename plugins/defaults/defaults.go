// Package defaults provides flags support for uconfig
package defaults

import "github.com/omeid/uconfig/flat"

const tag = "default"

// Defaults is an env variable plugin.
type Defaults interface {
	Visit(flat.Fields) error

	Parse() error
}

// New returns an EnvSet.
func New() Defaults {
	return &visitor{}
}

type visitor struct {
	fields flat.Fields
}

func (v *visitor) Visit(f flat.Fields) error {

	v.fields = f

	return nil
}

func (v *visitor) Parse() error {

	for _, f := range v.fields {
		value, ok := f.Tag(tag)
		if !ok {
			continue
		}

		f.Meta()[tag] = value
		err := f.Set(value)
		if err != nil {
			return err
		}
	}

	return nil
}
