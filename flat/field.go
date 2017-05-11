package flat

import (
	"reflect"
	"strconv"
)

var _ Field = (*field)(nil)

type field struct {
	name string
	meta map[string]string

	tag   reflect.StructTag
	field reflect.Value
}

// Used by standard library flag package.
func (f *field) IsBoolFlag() bool {
	return f.field.Kind() == reflect.Bool
}

func (f *field) Name() string {
	return f.name
}

func (f *field) Meta() map[string]string {
	return f.meta
}

func (f *field) Tag(key string) (string, bool) {
	return f.tag.Lookup(key)
}

func (f *field) String() string {
	return f.tag.Get("default")
}

func (f *field) Get() interface{} {
	// return f.ptr
	return nil
}

func (f *field) Set(value string) error {

	switch f.field.Kind() {
	case reflect.String:
		return f.setString(value)
	case reflect.Bool:
		return f.setBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.setInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.setUint(value)
	case reflect.Float32, reflect.Float64:
		return f.setFloat(value)

		// Soon case reflect.Map:
		// Soon case reflect.Slice:

		// Maybe case reflect.Func:
		// Maybe case reflect.Array:

		// Why? case reflect.Complex64:
		// Why? case reflect.Complex128:

		// Never case reflect.Chan:
		// Never case reflect.Interface:
		// Never case reflect.Ptr:
		// Never case reflect.Struct:
		// Never case reflect.UnsafePointer:
	}
	return nil
}

func (f *field) setString(value string) error {
	f.field.SetString(value)
	return nil
}

func (f *field) setBool(value string) error {
	v, err := strconv.ParseBool(value)
	f.field.SetBool(v)
	return err
}

func (f *field) setInt(value string) error {
	v, err := strconv.ParseInt(value, 0, 64)
	f.field.SetInt(v)
	return err
}

func (f *field) setUint(value string) error {
	v, err := strconv.ParseUint(value, 0, 64)
	f.field.SetUint(v)
	return err
}

func (f *field) setFloat(value string) error {
	v, err := strconv.ParseFloat(value, 64)
	f.field.SetFloat(v)
	return err
}
