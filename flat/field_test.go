package flat_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig/flat"
)

func TestField(t *testing.T) {

	type Config struct {
		First  string `test:"test-tag"`
		Second error
	}

	conf := &Config{First: "first"}
	fs, err := flat.View(conf)

	if err != nil {
		t.Fatal(err)
	}

	firstField := fs[0]

	name, _ := firstField.Name("")
	if name != "First" {
		t.Errorf("expected First but got %v", name)
	}

	tag, ok := firstField.Tag("test")
	if !ok {
		t.Error("expected test tag on firstField but not found")
	}

	if tag != "test-tag" {
		t.Errorf("expected tag test to be (test-tag) but got (%v)", tag)
	}

	meta1 := firstField.Meta()
	meta2 := firstField.Meta()

	meta1["test"] = "okay"

	if diff := cmp.Diff(meta1, meta2); diff != "" {
		t.Error(diff)
	}

	if def := firstField.Interface(); def != "first" {
		t.Errorf("expected Interface() to return default tag value (first) but got (%v)", def)
	}

	firstFieldPtr := firstField.Ptr().(*string)
	*firstFieldPtr = "first via pointer"

	if def := firstField.Interface(); def != "first via pointer" {
		t.Errorf("expected String() to return value set via pointer but got %v", def)
	}

}
