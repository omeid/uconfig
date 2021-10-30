package required

import (
	"fmt"
	"reflect"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

type ErrRequiredField struct {
	Field string
}

func (e ErrRequiredField) Error() string {
	return fmt.Sprintf("%s requires a value", e.Field)
}

const tag = "required"

func init() {
	plugins.RegisterTag(tag)
}

func New() plugins.Plugin {
	return &visitor{}
}

type visitor struct {
	fields flat.Fields
}

func (v *visitor) Parse() error {
	for _, field := range v.fields {
		value, ok := field.Tag(tag)

		if ok && value == "true" {
			fieldType := reflect.TypeOf(field.Get())

			if fieldType != nil && fieldType.Comparable() && field.Get() == reflect.Zero(fieldType).Interface() {
				return ErrRequiredField{Field: field.Name()}
			}
		}
	}

	return nil
}

func (v *visitor) Visit(fields flat.Fields) error {
	v.fields = fields

	for _, f := range v.fields {
		value, ok := f.Tag(tag)
		if ok {
			f.Meta()[tag] = value
		}
	}

	return nil
}
