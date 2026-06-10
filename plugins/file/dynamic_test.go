package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/omeid/uconfig/plugins/file"
)

func TestWorkspaceFindsInCWD(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".myapp", "config"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir: %v", err)
		}
	}()

	got := file.Workspace(".myapp/config").Resolve()
	want := filepath.Join(dir, ".myapp", "config")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceWalksUp(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".myapp", "config"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir: %v", err)
		}
	}()

	got := file.Workspace(".myapp/config").Resolve()
	want := filepath.Join(root, ".myapp", "config")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceReturnsClosest(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".myapp", "config"), []byte(`{"level":"root"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	child := filepath.Join(root, "project")
	if err := os.MkdirAll(filepath.Join(child, ".myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(child, ".myapp", "config"), []byte(`{"level":"child"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(child); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir: %v", err)
		}
	}()

	got := file.Workspace(".myapp/config").Resolve()
	want := filepath.Join(child, ".myapp", "config")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceNotFound(t *testing.T) {
	dir := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir: %v", err)
		}
	}()

	got := file.Workspace(".nonexistent/config").Resolve()
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestWorkspaceIsLazy(t *testing.T) {
	dir := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir: %v", err)
		}
	}()

	// Create the resolver before the file exists.
	resolve := file.Workspace(".myapp/config")

	if got := resolve.Resolve(); got != "" {
		t.Fatalf("expected empty before file exists, got %q", got)
	}

	// Now create the file.
	if err := os.MkdirAll(filepath.Join(dir, ".myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".myapp", "config"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(dir, ".myapp", "config")
	if got := resolve.Resolve(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWorkspaceDisplayName(t *testing.T) {
	p := file.Workspace(".myapp/config")
	if p.Name() != ".myapp/config" {
		t.Fatalf("expected name %q, got %q", ".myapp/config", p.Name())
	}
}

func TestAbsReturnsFixedPath(t *testing.T) {
	p := file.Absolute("/etc/app/config.json")
	if p.Resolve() != "/etc/app/config.json" {
		t.Fatalf("got %q, want %q", p.Resolve(), "/etc/app/config.json")
	}
	if p.Name() != "/etc/app/config.json" {
		t.Fatalf("name: got %q, want %q", p.Name(), "/etc/app/config.json")
	}
}

func TestRelativeResolvesAgainstCWD(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("chdir: %v", err)
		}
	}()

	p := file.Relative("config.json")
	got := p.Resolve()
	want := filepath.Join(dir, "config.json")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if p.Name() != "config.json" {
		t.Fatalf("name: got %q, want %q", p.Name(), "config.json")
	}
}
