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

func (f *field) Interface() any {
	return f.field.Interface()
}

func (f *field) Ptr() any {
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
	case reflect.Map:
		return f.setMap(value)

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
	setter := typeSetter(t.Elem())

	if setter == nil {
		return nil
	}

	values := strings.Split(value, ",")
	valuesLen := len(values)

	f.field.Set(reflect.MakeSlice(t, valuesLen, valuesLen))

	for i, value := range values {
		err := setter(f.field.Index(i), strings.TrimSpace(value))
		if err != nil {
			return err
		}
	}

	return nil
}

// setMap parses "key:value,key:value" into a map.
// Supports all types that typeSetter handles for both keys and values.
func (f *field) setMap(value string) error {
	t := f.field.Type()

	setKey := typeSetter(t.Key())
	if setKey == nil {
		return nil
	}

	setVal := typeSetter(t.Elem())
	if setVal == nil {
		return nil
	}

	m := reflect.MakeMap(t)

	entries := strings.Split(value, ",")
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		rawKey, rawVal, ok := strings.Cut(entry, ":")
		if !ok {
			continue
		}

		k := reflect.New(t.Key()).Elem()
		if err := setKey(k, strings.TrimSpace(rawKey)); err != nil {
			return err
		}

		v := reflect.New(t.Elem()).Elem()
		if err := setVal(v, strings.TrimSpace(rawVal)); err != nil {
			return err
		}

		m.SetMapIndex(k, v)
	}

	f.field.Set(m)
	return nil
}

func typeSetter(elem reflect.Type) func(reflect.Value, string) error {
	if elem.Implements(textUnmarshalerType) {
		return typeSetterUnmarshale
	}

	if reflect.PointerTo(elem).Implements(textUnmarshalerType) {
		return typeSetterPtrUnmarshale
	}

	switch elem.Kind() {

	case reflect.String:
		return typeSetterString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if elem.String() == "time.Duration" {
			return typeSetterDuration
		}

		return typeSetterInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return typeSetterUint

	case reflect.Float32, reflect.Float64:
		return typeSetterFloat
	}

	return nil
}

func typeSetterUnmarshale(f reflect.Value, value string) error {
	ptr := reflect.New(f.Type().Elem())

	ut := ptr.MethodByName("UnmarshalText")
	err := ut.Call([]reflect.Value{reflect.ValueOf([]byte(value))})[0]

	if !err.IsNil() {
		return err.Interface().(error)
	}

	f.Set(ptr)
	return nil
}

func typeSetterPtrUnmarshale(f reflect.Value, value string) error {
	ptr := reflect.New(f.Type())

	ut := ptr.MethodByName("UnmarshalText")
	err := ut.Call([]reflect.Value{reflect.ValueOf([]byte(value))})[0]

	if !err.IsNil() {
		return err.Interface().(error)
	}

	f.Set(reflect.Indirect(ptr))
	return nil
}

func typeSetterString(f reflect.Value, value string) error {
	f.SetString(value)
	return nil
}

func typeSetterDuration(f reflect.Value, value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	f.SetInt(int64(duration))
	return nil
}

func typeSetterInt(f reflect.Value, value string) error {
	v, err := strconv.ParseInt(value, 0, 64)
	f.SetInt(v)
	return err
}

func typeSetterUint(f reflect.Value, value string) error {
	v, err := strconv.ParseUint(value, 0, 64)
	f.SetUint(v)
	return err
}

func typeSetterFloat(f reflect.Value, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	f.SetFloat(v)
	return err
}
