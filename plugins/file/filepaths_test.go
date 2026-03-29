package file_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/file"
)

func TestFilePathsFromFilePlugins(t *testing.T) {
	ps := []plugins.Plugin{
		file.New("config.json", json.Unmarshal, file.Config{Optional: true}),
		file.New("/etc/app/config.json", json.Unmarshal, file.Config{Optional: true}),
	}

	paths := file.FilePaths(ps)
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "config.json" {
		t.Fatalf("expected config.json, got %s", paths[0])
	}
	if paths[1] != "/etc/app/config.json" {
		t.Fatalf("expected /etc/app/config.json, got %s", paths[1])
	}
}

func TestFilePathsIgnoresNonFilePlugins(t *testing.T) {
	ps := []plugins.Plugin{
		defaults.New(),
		file.New("config.json", json.Unmarshal, file.Config{Optional: true}),
	}

	paths := file.FilePaths(ps)
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d: %v", len(paths), paths)
	}
}

func TestFilePathsEmpty(t *testing.T) {
	paths := file.FilePaths(nil)
	if len(paths) != 0 {
		t.Fatalf("expected 0 paths, got %d", len(paths))
	}

	paths = file.FilePaths([]plugins.Plugin{defaults.New()})
	if len(paths) != 0 {
		t.Fatalf("expected 0 paths for non-file plugins, got %d", len(paths))
	}
}

func TestParseReReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	type Config struct {
		Name string `json:"name"`
	}

	// Write initial config.
	os.WriteFile(path, []byte(`{"name":"v1"}`), 0644)

	plug := file.New(path, json.Unmarshal, file.Config{})
	conf := &Config{}
	plug.(plugins.Walker).Walk(conf)

	// First parse.
	if err := plug.Parse(); err != nil {
		t.Fatal(err)
	}
	if conf.Name != "v1" {
		t.Fatalf("expected v1, got %s", conf.Name)
	}

	// Change the file.
	os.WriteFile(path, []byte(`{"name":"v2"}`), 0644)

	// Re-walk with fresh struct and re-parse -- should read new content.
	conf2 := &Config{}
	plug.(plugins.Walker).Walk(conf2)
	if err := plug.Parse(); err != nil {
		t.Fatal(err)
	}
	if conf2.Name != "v2" {
		t.Fatalf("expected v2, got %s", conf2.Name)
	}
}

func TestParseOptionalMissing(t *testing.T) {
	plug := file.New("/nonexistent/config.json", json.Unmarshal, file.Config{Optional: true})
	conf := &struct{}{}
	plug.(plugins.Walker).Walk(conf)

	// Should not error.
	if err := plug.Parse(); err != nil {
		t.Fatalf("optional missing file should not error: %v", err)
	}
}

func TestParseRequiredMissing(t *testing.T) {
	plug := file.New("/nonexistent/config.json", json.Unmarshal, file.Config{Optional: false})
	conf := &struct{}{}

	// Walk should catch missing required file.
	if err := plug.(plugins.Walker).Walk(conf); err == nil {
		t.Fatal("required missing file should error in Walk")
	}
}

func TestParseMultipleTimes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	type Config struct {
		Value int `json:"value"`
	}

	os.WriteFile(path, []byte(`{"value":1}`), 0644)

	plug := file.New(path, json.Unmarshal, file.Config{})

	// Parse 3 times, changing file each time.
	for i := 1; i <= 3; i++ {
		os.WriteFile(path, []byte(fmt.Sprintf(`{"value":%d}`, i)), 0644)
		conf := &Config{}
		plug.(plugins.Walker).Walk(conf)
		if err := plug.Parse(); err != nil {
			t.Fatal(err)
		}
		if conf.Value != i {
			t.Fatalf("parse %d: expected %d, got %d", i, i, conf.Value)
		}
	}
}
