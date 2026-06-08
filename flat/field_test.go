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

	// Non-collection fields should return an error.
	err = firstField.Append("key", func(v any) error { return nil })
	if err == nil {
		t.Error("expected Append() to return error for a string field")
	}
}

func TestAppendMap(t *testing.T) {
	type Item struct {
		Name string
	}
	type Config struct {
		Items map[string]Item
	}

	conf := &Config{}
	fs, err := flat.View(conf)
	if err != nil {
		t.Fatal(err)
	}

	f := fs[0]

	err = f.Append("a.json", func(v any) error {
		item := v.(*Item)
		item.Name = "alpha"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	err = f.Append("b.json", func(v any) error {
		item := v.(*Item)
		item.Name = "beta"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]Item{
		"a.json": {Name: "alpha"},
		"b.json": {Name: "beta"},
	}
	if diff := cmp.Diff(expected, conf.Items); diff != "" {
		t.Error(diff)
	}
}

func TestAppendSlice(t *testing.T) {
	type Item struct {
		Name string
	}
	type Config struct {
		Items []Item
	}

	conf := &Config{}
	fs, err := flat.View(conf)
	if err != nil {
		t.Fatal(err)
	}

	f := fs[0]

	err = f.Append("ignored-key", func(v any) error {
		item := v.(*Item)
		item.Name = "a"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	err = f.Append("another-ignored-key", func(v any) error {
		item := v.(*Item)
		item.Name = "b"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := []Item{
		{Name: "a"},
		{Name: "b"},
	}
	if diff := cmp.Diff(expected, conf.Items); diff != "" {
		t.Error(diff)
	}
}

func TestDeleteMap(t *testing.T) {
	type Item struct {
		Name string
	}
	type Config struct {
		Items map[string]Item
	}

	conf := &Config{
		Items: map[string]Item{
			"a.json": {Name: "alpha"},
			"b.json": {Name: "beta"},
		},
	}
	fs, err := flat.View(conf)
	if err != nil {
		t.Fatal(err)
	}

	f := fs[0]

	err = f.Delete("a.json")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]Item{
		"b.json": {Name: "beta"},
	}
	if diff := cmp.Diff(expected, conf.Items); diff != "" {
		t.Error(diff)
	}
}

func TestDeleteSlice(t *testing.T) {
	type Item struct {
		Name string
	}
	type Config struct {
		Items []Item
	}

	conf := &Config{
		Items: []Item{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
	}
	fs, err := flat.View(conf)
	if err != nil {
		t.Fatal(err)
	}

	f := fs[0]

	// Delete index 1 ("b")
	err = f.Delete("1")
	if err != nil {
		t.Fatal(err)
	}

	expected := []Item{
		{Name: "a"},
		{Name: "c"},
	}
	if diff := cmp.Diff(expected, conf.Items); diff != "" {
		t.Error(diff)
	}

	// Delete invalid index
	err = f.Delete("2")
	if err == nil {
		t.Error("expected error for out of bounds index")
	}

	// Delete invalid integer
	err = f.Delete("not-a-number")
	if err == nil {
		t.Error("expected error for non-integer slice key")
	}
}

func TestAppendMapError(t *testing.T) {
	type Config struct {
		Items map[int]string
	}

	conf := &Config{}
	fs, err := flat.View(conf)
	if err != nil {
		t.Fatal(err)
	}

	err = fs[0].Append("1", func(v any) error { return nil })
	if err == nil {
		t.Error("expected Append() to return error for map[int]string")
	}
}
