package file_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
)

func TestRequireOne(t *testing.T) {
	os.Args = os.Args[:1]
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "config1.json")
	file2 := filepath.Join(tempDir, "config2.json")

	validJSON := `{"Command":"run-exclusive", "GoHard":true}`

	// Test Case 1: No files exist -> ErrNoFiles
	files := file.RequireOne{
		{file.Absolute(file1), json.Unmarshal},
		{file.Absolute(file2), json.Unmarshal},
	}

	conf1 := uconfig.Classic[f.Config](files)
	_, err := conf1.Parse()
	if err == nil || err.Error() != file.ErrNoFiles.Error() {
		t.Fatalf("expected ErrNoFiles, got: %v", err)
	}

	// Test Case 2: Exactly one file exists -> Success
	err = os.WriteFile(file1, []byte(validJSON), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	conf2 := uconfig.Classic[f.Config](files)
	val, err := conf2.Parse()
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}

	if val.Command != "run-exclusive" || val.GoHard != true {
		t.Errorf("parsed value incorrect: %+v", val)
	}

	// Test Case 3: Multiple files exist -> ErrMultipleFiles
	err = os.WriteFile(file2, []byte(validJSON), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	conf3 := uconfig.Classic[f.Config](files)
	_, err = conf3.Parse()
	if err == nil || !errors.Is(err, file.ErrMultipleFiles) {
		t.Fatalf("expected ErrMultipleFiles, got: %v", err)
	}
}
