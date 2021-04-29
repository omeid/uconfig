package flat

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/omeid/uconfig/internal/f"
)

func TestFlattenNested(t *testing.T) {

	conf := f.Config{}
	fs, err := View(&conf)

	if err != nil {
		t.Fatal(err)
	}

	fields := map[string]bool{
		"GoHard":               false,
		"Version":              false,
		"Redis.Host":           false,
		"Redis.Port":           false,
		"Rethink.Host.Address": false,
		"Rethink.Host.Port":    false,
		"Rethink.Db":           false,
	}

	// for _, fs := range fs {
	// 	t.Log(" - ", fs.Name())
	// }

	for _, fs := range fs {
		name := fs.Name()
		_, ok := fields[name]
		if !ok {
			t.Fatalf("Unexpected Field: %v", name)
		}

		fields[name] = true
	}

	for name, seen := range fields {
		if !seen {
			t.Fatalf("Expected field missing: %v", name)
		}
	}
}

func TestFlattenTypes(t *testing.T) {
	values := map[string]string{
		"String":  "string",
		"Bool":    "true",
		"Int":     "1",
		"Int8":    "2",
		"Int16":   "4",
		"Int32":   "8",
		"Int64":   "16",
		"Uint":    "32",
		"Uint8":   "64",
		"Uint16":  "128",
		"Uint32":  "256",
		"Uint64":  "512",
		"Float32": "1.1",
		"Float64": "1.2",
		// "Duration":        "5 second",
		// "MapStringString": "a:aval,bbval",
		// "MapStringInt":    "one:1,two:2",
		"SliceString": "hello,world",
		// "SliceInt":        "1,2,3",
	}

	expect := f.Types{
		String:  "string",
		Bool:    true,
		Int:     1,
		Int8:    2,
		Int16:   4,
		Int32:   8,
		Int64:   16,
		Uint:    32,
		Uint8:   64,
		Uint16:  128,
		Uint32:  256,
		Uint64:  512,
		Float32: 1.1,
		Float64: 1.2,

		// Duration:        time.Second * 5,
		// MapStringString: map[string]string{"a": "aval", "b": "bval"},
		// MapStringInt:    map[string]int{"one": 1, "two": 2},
		SliceString: []string{"hello", "world"},
		// SliceInt:        []int{1, 2, 3},
	}

	value := f.Types{}

	_ = values

	fs, err := View(&value)

	if err != nil {
		t.Fatal(err)
	}

	for _, field := range fs {
		name := field.Name()
		value, ok := values[name]
		if !ok {
			t.Fatalf("Missing value for %v", name)
		}

		err := field.Set(value)

		if err != nil {
			t.Fatalf("Field: %v, Error: %v", name, err)
		}
	}

	if diff := deep.Equal(expect, value); diff != nil {
		t.Error(diff)
	}

}
