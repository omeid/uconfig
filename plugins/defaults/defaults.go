// Package defaults provides flags support for uconfig
package defaults

import (
	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

const tag = "default"

func init() {
	plugins.RegisterTag(tag)
}

// New returns an EnvSet.
func New() plugins.Plugin {
	return &visitor{}
}

type visitor struct {
	fields flat.Fields
}

func (v *visitor) Visit(f flat.Fields) error {
	v.fields = f

	for _, f := range v.fields {
		value, ok := f.Tag(tag)
		if !ok {
			continue
		}
		f.Meta()[tag] = value
	}
	return nil
}

func (v *visitor) Parse() error {
	for _, f := range v.fields {
		value, ok := f.Meta()[tag]
		if !ok {
			continue
		}
		err := f.Set(value)
		if err != nil {
			return err
		}
	}

	return nil
}
