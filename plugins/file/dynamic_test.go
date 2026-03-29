package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/omeid/uconfig/plugins/file"
)

func TestWorkspaceFindsInCWD(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".myapp"), 0755)
	os.WriteFile(filepath.Join(dir, ".myapp", "config"), []byte("{}"), 0644)

	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	got := file.Workspace(".myapp/config").Resolve()
	want := filepath.Join(dir, ".myapp", "config")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceWalksUp(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".myapp"), 0755)
	os.WriteFile(filepath.Join(root, ".myapp", "config"), []byte("{}"), 0644)

	nested := filepath.Join(root, "a", "b", "c")
	os.MkdirAll(nested, 0755)

	orig, _ := os.Getwd()
	os.Chdir(nested)
	defer os.Chdir(orig)

	got := file.Workspace(".myapp/config").Resolve()
	want := filepath.Join(root, ".myapp", "config")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceReturnsClosest(t *testing.T) {
	root := t.TempDir()

	os.MkdirAll(filepath.Join(root, ".myapp"), 0755)
	os.WriteFile(filepath.Join(root, ".myapp", "config"), []byte(`{"level":"root"}`), 0644)

	child := filepath.Join(root, "project")
	os.MkdirAll(filepath.Join(child, ".myapp"), 0755)
	os.WriteFile(filepath.Join(child, ".myapp", "config"), []byte(`{"level":"child"}`), 0644)

	orig, _ := os.Getwd()
	os.Chdir(child)
	defer os.Chdir(orig)

	got := file.Workspace(".myapp/config").Resolve()
	want := filepath.Join(child, ".myapp", "config")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceNotFound(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	got := file.Workspace(".nonexistent/config").Resolve()
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestWorkspaceIsLazy(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// Create the resolver before the file exists.
	resolve := file.Workspace(".myapp/config")

	if got := resolve.Resolve(); got != "" {
		t.Fatalf("expected empty before file exists, got %q", got)
	}

	// Now create the file.
	os.MkdirAll(filepath.Join(dir, ".myapp"), 0755)
	os.WriteFile(filepath.Join(dir, ".myapp", "config"), []byte("{}"), 0644)

	want := filepath.Join(dir, ".myapp", "config")
	if got := resolve.Resolve(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceDisplayName(t *testing.T) {
	p := file.Workspace(".myapp/config")
	if p.Name != "workspace: .myapp/config" {
		t.Fatalf("expected name %q, got %q", "workspace: .myapp/config", p.Name)
	}
}

func TestAbsReturnsFixedPath(t *testing.T) {
	p := file.Absolute("/etc/app/config.json")
	if p.Resolve() != "/etc/app/config.json" {
		t.Fatalf("got %q, want %q", p.Resolve(), "/etc/app/config.json")
	}
	if p.Name != "absolute:  /etc/app/config.json" {
		t.Fatalf("name: got %q, want %q", p.Name, "absolute:  /etc/app/config.json")
	}
}

func TestRelativeResolvesAgainstCWD(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	p := file.Relative("config.json")
	got := p.Resolve()
	want := filepath.Join(dir, "config.json")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if p.Name != "relative:  config.json" {
		t.Fatalf("name: got %q, want %q", p.Name, "relative:  config.json")
	}
}
