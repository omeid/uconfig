// Package defaults provides default value support for uconfig.
package defaults

import (
	"slices"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

const tag = "default"

func init() {
	plugins.RegisterTag(tag)
}

// New returns a defaults plugin.
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

func SortDefaultsFirst(ps []plugins.Plugin) []plugins.Plugin {
	if len(ps) < 2 {
		return ps
	}

	idx := -1
	for i, p := range ps {
		if _, ok := p.(*visitor); ok {
			idx = i
			break
		}
	}

	if idx < 1 {
		// no defaults or already at zero (eg. Classic).
		return ps
	}

	return slices.Concat(
		ps[idx:idx+1],
		ps[:idx],
		ps[idx+1:],
	)
}
