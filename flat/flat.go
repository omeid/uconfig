// Package flat provides a flat view of an arbitrary nested structs.
package flat

import (
	"errors"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ErrUnexpectedType is returned when flatten sees an unsupported type.
var ErrUnexpectedType = errors.New("unexpected type, expecting a pointer to struct")

// Fields is a slice of Field.
type Fields []Field

// Field describe an interface to our flat structs fields.
type Field interface {
	// Name returns the name for a given tag, if any
	// and also whatever the returned name is "explicit" by
	// the user or plugins are allowed to rewrite it.
	Name(tag string) (string, bool)

	Tag(key string) (string, bool)

	Meta() map[string]string

	Interface() any
	Set(string) error

	// returns the Ptr to this value.
	// It is used by complex decoders like uconfig-cue.
	Ptr() any

	// Append adds a new entry to the field's collection.
	// It allocates a new zero value of the element type, calls fn
	// with a pointer to it for population (e.g. via unmarshal), and
	// then inserts the result into the collection.
	// For map fields, key is used as the map key.
	// For slice fields, key is ignored and the value is appended.
	// Returns an error if the field is not a collection type or if
	// fn returns an error.
	Append(key string, fn func(any) error) error

	// Delete removes an entry from the field's collection.
	// For map fields, key is the map key to remove.
	// For slice fields, key is the index to remove (parsed as an integer).
	// Returns an error if the field is not a collection type or if
	// the index is invalid.
	Delete(key string) error

	// Folded returns true if the field is a collection (map or slice)
	// whose elements cannot be parsed from a single string.
	// Such fields can only be populated via Append.
	Folded() bool
}

var caser = cases.Title(language.Und, cases.NoLower)

// View provides a flat view of the provided structs an array of fields.
// sub-struct fields are prefixed with the struct key (not type) followed by a dot,
// this is repeated for each nested level.
func View(s any) (Fields, error) {
	rs, err := unwrap(s)
	if err != nil {
		return nil, err
	}

	return walkStruct("", rs)
}

func walkStruct(prefix string, rs reflect.Value) ([]Field, error) {
	prefix = caser.String(prefix)

	fields := []Field{}

	ts := rs.Type()
	for i := 0; i < rs.NumField(); i++ {

		fv := rs.Field(i)
		ft := ts.Field(i)

		switch fv.Kind() {

		case reflect.Struct:
			structPrefix := prefix
			if !ft.Anonymous {
				// Unless it is anonymous struct, append the field name to the prefix.
				if structPrefix == "" {
					structPrefix = ft.Name
				} else {
					structPrefix = structPrefix + "." + ft.Name
				}
			}
			fs, err := walkStruct(structPrefix, fv)
			if err != nil {
				return nil, err
			}
			fields = append(fields, fs...)
		default:

			fieldName := ft.Name

			// unless it is override
			if name, ok := ft.Tag.Lookup("uconfig"); ok && name != "" {
				fieldName = name
			}

			fields = append(fields, &field{
				name:   fieldName,
				prefix: prefix,
				meta:   make(map[string]string, 5),
				tag:    ft.Tag,
				field:  fv,
			})
		}
	}

	return fields, nil
}

func unwrap(s any) (reflect.Value, error) {
	rs := reflect.ValueOf(s)

	if k := rs.Kind(); k != reflect.Pointer {
		return rs, ErrUnexpectedType
	}

	rs = reflect.Indirect(rs)

	if rs.Kind() == reflect.Interface {
		rs = rs.Elem()
	}

	rs = reflect.Indirect(rs)

	if rs.Kind() != reflect.Struct {
		return rs, ErrUnexpectedType
	}

	return rs, nil
}
