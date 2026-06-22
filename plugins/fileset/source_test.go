package fileset

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/omeid/uconfig/plugins/file"
)

func TestSourceAbsolute(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "uconfig-test-source-abs-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	files := []string{"a.json", "b.json", "c.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tempDir, f), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	pattern := filepath.Join(tempDir, "*.json")
	set := Absolute(pattern)

	if set.Name() != pattern {
		t.Errorf("Expected name %q, got %q", pattern, set.Name())
	}
	if set.Kind() != file.KindAbsolute {
		t.Errorf("Expected kind %v, got %v", file.KindAbsolute, set.Kind())
	}

	sources, err := set.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if len(sources) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(sources))
	}

	var resolved []string
	for _, s := range sources {
		resolved = append(resolved, s.Resolve())
	}
	sort.Strings(resolved)

	expected := []string{
		filepath.Join(tempDir, "a.json"),
		filepath.Join(tempDir, "b.json"),
	}

	if !reflect.DeepEqual(resolved, expected) {
		t.Errorf("Expected %v, got %v", expected, resolved)
	}
}

func TestSourceRelative(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "uconfig-test-source-rel-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir("apps.d", 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{"apps.d/a.json", "apps.d/b.json", "apps.d/c.txt"}
	for _, f := range files {
		if err := os.WriteFile(f, []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	pattern := "apps.d/*.json"
	set := Relative(pattern)

	if set.Name() != pattern {
		t.Errorf("Expected name %q, got %q", pattern, set.Name())
	}
	if set.Kind() != file.KindRelative {
		t.Errorf("Expected kind %v, got %v", file.KindRelative, set.Kind())
	}

	sources, err := set.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if len(sources) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(sources))
	}

	var resolved []string
	for _, s := range sources {
		resolved = append(resolved, s.Resolve())
	}
	sort.Strings(resolved)

	expected := []string{
		filepath.Join(tempDir, "apps.d", "a.json"),
		filepath.Join(tempDir, "apps.d", "b.json"),
	}

	if !reflect.DeepEqual(resolved, expected) {
		t.Errorf("Expected %v, got %v", expected, resolved)
	}
}

func TestSourceWorkspace(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "uconfig-test-source-workspace-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	workspaceRoot := tempDir
	configDir := filepath.Join(workspaceRoot, ".myapp")
	nestedDir := filepath.Join(workspaceRoot, "nested", "dir")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{"a.json", "b.json", "c.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(configDir, f), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.Chdir(nestedDir); err != nil {
		t.Fatal(err)
	}

	pattern := ".myapp/*.json"
	set := Workspace(pattern)

	if set.Name() != pattern {
		t.Errorf("Expected name %q, got %q", pattern, set.Name())
	}
	if set.Kind() != file.KindWorkspace {
		t.Errorf("Expected kind %v, got %v", file.KindWorkspace, set.Kind())
	}

	sources, err := set.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if len(sources) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(sources))
	}

	var resolved []string
	for _, s := range sources {
		resolved = append(resolved, s.Resolve())
	}
	sort.Strings(resolved)

	expected := []string{
		filepath.Join(configDir, "a.json"),
		filepath.Join(configDir, "b.json"),
	}

	if !reflect.DeepEqual(resolved, expected) {
		t.Errorf("Expected %v, got %v", expected, resolved)
	}
}

func TestSourceWorkspaceNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "uconfig-test-source-workspace-notfound-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	set := Workspace(".myapp/*.json")
	sources, err := set.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if len(sources) != 0 {
		t.Errorf("Expected 0 sources, got %d", len(sources))
	}
}
