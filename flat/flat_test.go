package flat_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/internal/f"
)

func TestFlattenNested(t *testing.T) {

	conf := f.Config{}
	fs, err := flat.View(&conf)

	if err != nil {
		t.Fatal(err)
	}

	fields := map[string]bool{
		"GoHard":               false,
		"Version":              false,
		"Redis.Address":        false,
		"Redis.Port":           false,
		"Rethink.Host.Address": false,
		"Rethink.Host.Port":    false,
		"Rethink.Db":           false,
		"Rethink.Password":     false,
	}

	for _, fs := range fs {
		name, explicit := fs.Name("")
		t.Log(" - ", name, explicit)
	}

	for _, fs := range fs {
		name, _ := fs.Name("")
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
		"String":   "string",
		"Bool":     "true",
		"Int":      "1",
		"Int8":     "2",
		"Int16":    "4",
		"Int32":    "8",
		"Int64":    "16",
		"Uint":     "32",
		"Uint8":    "64",
		"Uint16":   "128",
		"Uint32":   "256",
		"Uint64":   "512",
		"Float32":  "1.1",
		"Float64":  "1.2",
		"Duration": "5s",

		// "MapStringString": "a:aval,bbval",
		// "MapStringInt":    "one:1,two:2",

		"SliceString":   "hello,world",
		"SliceInt":      "1, 2,3",
		"SliceInt32":    "1,2, 3 ",
		"SliceUint":     "1, 2,3",
		"SliceFloat32":  "1.2,3.4, 5.6",
		"SliceDuration": "5s, 1h",

		"SliceTextUnmarshaler":    "a.b.c",
		"SliceElemUnmarshaler":    "north,east,south,west",
		"SliceElemPtrUnmarshaler": "north,east,south,west",
	}

	expect := f.Types{
		String:   "string",
		Bool:     true,
		Int:      1,
		Int8:     2,
		Int16:    4,
		Int32:    8,
		Int64:    16,
		Uint:     32,
		Uint8:    64,
		Uint16:   128,
		Uint32:   256,
		Uint64:   512,
		Float32:  1.1,
		Float64:  1.2,
		Duration: time.Duration(5 * time.Second),

		// Duration:        time.Second * 5,
		// MapStringString: map[string]string{"a": "aval", "b": "bval"},
		// MapStringInt:    map[string]int{"one": 1, "two": 2},
		SliceString:   []string{"hello", "world"},
		SliceInt:      []int{1, 2, 3},
		SliceInt32:    []int{1, 2, 3},
		SliceUint:     []uint{1, 2, 3},
		SliceFloat32:  []float32{1.2, 3.4, 5.6},
		SliceDuration: []time.Duration{5 * time.Second, 1 * time.Hour},

		SliceTextUnmarshaler: &f.TextUnmarshalerStringSlice{"a", "b", "c"},
		SliceElemUnmarshaler: f.ElemUnmarshalerSlice{0, 90, 180, 270},

		SliceElemPtrUnmarshaler: f.ElemPtrUnmarshalerSlice{
			f.NewReadableDirection(0),
			f.NewReadableDirection(90),
			f.NewReadableDirection(180),
			f.NewReadableDirection(270),
		},
	}

	value := f.Types{}

	_ = values

	fs, err := flat.View(&value)

	if err != nil {
		t.Fatal(err)
	}

	for _, field := range fs {
		name, _ := field.Name("")

		value, ok := values[name]
		if !ok {
			t.Fatalf("Missing value for %v", name)
		}

		fmt.Printf("Mapping field %s: %s\n", name, value)
		err := field.Set(value)

		if err != nil {
			t.Fatalf("Field: %v, Error: %v", name, err)
		}
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}
