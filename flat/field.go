package flat

import (
	"encoding"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var _ Field = (*field)(nil)

type field struct {
	name   string
	prefix string

	meta map[string]string

	tag   reflect.StructTag
	field reflect.Value
}

func (f *field) getName(tag string) (string, bool) {

	name, explicit := f.Tag(tag)

	name, _, _ = strings.Cut(name, ",")

	if name == "" || name == "." {
		name = f.name
		// explicit here means what it is an explicit name or should
		// be prefixed.
		explicit = false
	}

	if name[0] == '.' {
		name = name[1:]
	}

	return name, explicit
}

func (f *field) Name(tag string) (string, bool) {
	name, explicit := f.getName(tag)

	if f.prefix == "" || explicit {
		return name, explicit
	}

	return f.prefix + "." + name, explicit
}

func (f *field) Meta() map[string]string {
	return f.meta
}

func (f *field) Tag(key string) (string, bool) {
	if key == "" {
		return "", false
	}
	return f.tag.Lookup(key)
}

func (f *field) Interface() interface{} {
	return f.field.Interface()
}

func (f *field) Ptr() interface{} {

	kind := f.field.Kind()

	if kind == reflect.Pointer || kind == reflect.Slice || kind == reflect.Interface {
		return f.field.Interface()
	}

	return f.field.Addr().Interface()
}

var textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func (f *field) Set(value string) error {

	t := f.field.Type()

	if t.Implements(textUnmarshalerType) {
		return f.setUnmarshale([]byte(value))
	}

	switch f.field.Kind() {
	case reflect.String:
		return f.setString(value)
	case reflect.Bool:
		return f.setBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.String() == "time.Duration" {
			return f.setDuration(value)
		}
		return f.setInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.setUint(value)
	case reflect.Float32, reflect.Float64:
		return f.setFloat(value)
	case reflect.Slice:
		return f.setSlice(value)

		// Soon case reflect.Map:

		// Maybe case reflect.Array:

		// Why? case reflect.Complex64:
		// Why? case reflect.Complex128:

		// Never case reflect.Func:
		// Never case reflect.Chan:
		// Never case reflect.Interface:
		// Never case reflect.Ptr:
		// Never case reflect.Struct:
		// Never case reflect.UnsafePointer:
	}
	return nil
}

func (f *field) setUnmarshale(value []byte) error {

	if f.field.IsNil() {
		f.field.Set(reflect.New(f.field.Type().Elem()))
	}

	ut := f.field.MethodByName("UnmarshalText")

	err := ut.Call([]reflect.Value{reflect.ValueOf(value)})[0]

	if err.IsNil() {
		return nil
	}
	return err.Interface().(error)
}

func (f *field) setDuration(value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	f.field.SetInt(int64(duration))
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

func (f *field) setSlice(value string) error {

	t := f.field.Type()
	setSliceElem := setSliceElem(t.Elem())

	if setSliceElem == nil {
		return nil
	}

	values := strings.Split(value, ",")
	valuesLen := len(values)

	f.field.Set(reflect.MakeSlice(t, valuesLen, valuesLen))

	for i, value := range values {
		err := setSliceElem(f.field.Index(i), strings.TrimSpace(value))
		if err != nil {
			return err
		}
	}

	return nil
}

func setSliceElem(elem reflect.Type) func(reflect.Value, string) error {

	if elem.Implements(textUnmarshalerType) {
		return setSliceElemUnmarshale
	}

	if reflect.PointerTo(elem).Implements(textUnmarshalerType) {
		return setSliceElemPtrUnmarshale
	}

	switch elem.Kind() {

	case reflect.String:
		return setSliceElemString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if elem.String() == "time.Duration" {
			return setSliceElemDuration
		}

		return setSliceElemInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setSliceElemUint

	case reflect.Float32, reflect.Float64:
		return setSliceElemFloat
	}

	return nil
}

func setSliceElemUnmarshale(f reflect.Value, value string) error {
	ptr := reflect.New(f.Type().Elem())

	ut := ptr.MethodByName("UnmarshalText")
	err := ut.Call([]reflect.Value{reflect.ValueOf([]byte(value))})[0]

	if !err.IsNil() {
		return err.Interface().(error)
	}

	f.Set(ptr)
	return nil
}

func setSliceElemPtrUnmarshale(f reflect.Value, value string) error {
	ptr := reflect.New(f.Type())

	ut := ptr.MethodByName("UnmarshalText")
	err := ut.Call([]reflect.Value{reflect.ValueOf([]byte(value))})[0]

	if !err.IsNil() {
		return err.Interface().(error)
	}

	f.Set(reflect.Indirect(ptr))
	return nil
}

func setSliceElemString(f reflect.Value, value string) error {
	f.SetString(value)
	return nil
}

func setSliceElemDuration(f reflect.Value, value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	f.SetInt(int64(duration))
	return nil
}

func setSliceElemInt(f reflect.Value, value string) error {
	v, err := strconv.ParseInt(value, 0, 64)
	f.SetInt(v)
	return err
}

func setSliceElemUint(f reflect.Value, value string) error {
	v, err := strconv.ParseUint(value, 0, 64)
	f.SetUint(v)
	return err
}

func setSliceElemFloat(f reflect.Value, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	f.SetFloat(v)
	return err
}
