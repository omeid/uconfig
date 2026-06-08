package fileset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins"
)

type App struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type TestConfig struct {
	Apps []App `fileset:"my-apps"`
}

type InvalidConfig struct {
	Apps string `fileset:"my-apps"`
}

type DuplicateConfig struct {
	Apps1 []App `fileset:"my-apps"`
	Apps2 []App `fileset:"my-apps"`
}

type MapConfig struct {
	Apps map[string]App `fileset:"my-apps"`
}

type MapSliceConfig struct {
	Apps map[string][]App `fileset:"my-apps"`
}

func TestFilesetPlugin(t *testing.T) {
	// 1. Create a temporary directory
	tempDir, err := os.MkdirTemp("", "fileset_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 2. Write test JSON files (one object per file for []App)
	app1Content := `{"name": "auth", "port": 8080}`
	app2Content := `{"name": "users", "port": 8081}`
	app3Content := `{"name": "gateway", "port": 9000}`

	err = os.WriteFile(filepath.Join(tempDir, "app1.json"), []byte(app1Content), 0644)
	if err != nil {
		t.Fatalf("failed to write app1.json: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "app2.json"), []byte(app2Content), 0644)
	if err != nil {
		t.Fatalf("failed to write app2.json: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "app3.json"), []byte(app3Content), 0644)
	if err != nil {
		t.Fatalf("failed to write app3.json: %v", err)
	}

	// 3. Setup fileset plugin
	pattern := filepath.Join(tempDir, "*.json")
	p := New("my-apps", Absolute(pattern), json.Unmarshal)

	// 4. Load with uconfig
	conf := uconfig.New[TestConfig](p)
	parsed, err := conf.Parse()
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	// 5. Verify parsed results
	expected := []App{
		{Name: "auth", Port: 8080},
		{Name: "users", Port: 8081},
		{Name: "gateway", Port: 9000},
	}

	// Order returned by filepath.Glob is alphabetical, so app1.json then app2.json is guaranteed.
	if !reflect.DeepEqual(parsed.Apps, expected) {
		t.Errorf("expected parsed apps to be %v, got %v", expected, parsed.Apps)
	}

	// 6. Verify SourcePaths method
	sp, ok := p.(plugins.SourcePaths)
	if !ok {
		t.Fatalf("plugin does not implement plugins.SourcePaths")
	}
	srcs := sp.SourcePaths()
	if len(srcs) != 1 {
		t.Fatalf("expected 1 source, got %d", len(srcs))
	}
	expectedDir, _ := filepath.Abs(tempDir)
	if srcs[0].Path != expectedDir {
		t.Errorf("expected directory %q, got %q", expectedDir, srcs[0].Path)
	}
}

func TestUsage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "uconfig-test-usage-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file1 := filepath.Join(tempDir, "config1.json")
	file2 := filepath.Join(tempDir, "config2.json")
	
	if err := os.WriteFile(file1, []byte(`{"A": "1"}`), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := os.WriteFile(file2, []byte(`{"A": "2"}`), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	p := New("my apps", Path{
		Name: "test: " + filepath.Join(tempDir, "*.json"),
		Kind: "absolute",
		Resolve: func() ([]string, error) {
			return []string{file1, file2}, nil
		},
	}, json.Unmarshal)

	u, ok := p.(plugins.Usage)
	if !ok {
		t.Fatalf("plugin does not implement plugins.Usage")
	}

	// Test Usage(".")
	header, content := u.Usage(".")
	expectedHeader := "Files"
	if header != expectedHeader {
		t.Errorf("expected header %q, got %q", expectedHeader, header)
	}
	expectedContent := "test: " + filepath.Join(tempDir, "*.json") + "\n"
	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}

	// Test Usage("my apps")
	header, content = u.Usage("my apps")
	expectedHeader = "Fileset"
	if header != expectedHeader {
		t.Errorf("expected header %q, got %q", expectedHeader, header)
	}
	expectedContent = "absolute: test: " + filepath.Join(tempDir, "*.json") + "\t\n"
	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestRelativePath(t *testing.T) {
	// Create a temporary directory in the current working directory to test relative path
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}

	tempDirName := "test_temp_relative_dir"
	tempDir := filepath.Join(wd, tempDirName)
	err = os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Fatalf("failed to create relative temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	appContent := `{"name": "test", "port": 1234}`
	err = os.WriteFile(filepath.Join(tempDir, "app.json"), []byte(appContent), 0644)
	if err != nil {
		t.Fatalf("failed to write app.json: %v", err)
	}

	p := New("my-apps", Relative(filepath.Join(tempDirName, "*.json")), json.Unmarshal)
	conf := uconfig.New[TestConfig](p)
	parsed, err := conf.Parse()
	if err != nil {
		t.Fatalf("failed to parse config with relative path: %v", err)
	}

	expected := []App{{Name: "test", Port: 1234}}
	if !reflect.DeepEqual(parsed.Apps, expected) {
		t.Errorf("expected parsed apps to be %v, got %v", expected, parsed.Apps)
	}
}

func TestInvalidConfigTypes(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "fileset_invalid_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = os.WriteFile(filepath.Join(tempDir, "app.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("failed to write app.json: %v", err)
	}

	pattern := filepath.Join(tempDir, "*.json")
	p := New("my-apps", Absolute(pattern), json.Unmarshal)
	conf := uconfig.New[InvalidConfig](p)
	_, err = conf.Parse()
	if err == nil {
		t.Error("expected error when target field is not a slice or map, but got nil")
	}
}

func TestDuplicateConfigFields(t *testing.T) {
	p := New("my-apps", Absolute("*.json"), json.Unmarshal)
	conf := uconfig.New[DuplicateConfig](p)
	_, err := conf.Parse()
	if err == nil {
		t.Error("expected error when multiple fields have the same fileset tag, but got nil")
	}
}

func TestNoMatchingField(t *testing.T) {
	p := New("non-existent-tag", Absolute("*.json"), json.Unmarshal)
	conf := uconfig.New[TestConfig](p)
	_, err := conf.Parse()
	if err == nil {
		t.Error("expected error when no fields match the fileset tag, but got nil")
	}
}

func TestMapConfig(t *testing.T) {
	// 1. Create a temporary directory
	tempDir, err := os.MkdirTemp("", "fileset_map_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 2. Write test JSON files (single object per file for map[string]App)
	app1Content := `{"name": "auth", "port": 8080}`
	app2Content := `{"name": "gateway", "port": 9000}`

	err = os.WriteFile(filepath.Join(tempDir, "auth.json"), []byte(app1Content), 0644)
	if err != nil {
		t.Fatalf("failed to write auth.json: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "gateway.json"), []byte(app2Content), 0644)
	if err != nil {
		t.Fatalf("failed to write gateway.json: %v", err)
	}

	// 3. Setup fileset plugin
	pattern := filepath.Join(tempDir, "*.json")
	p := New("my-apps", Absolute(pattern), json.Unmarshal)

	// 4. Load with uconfig
	conf := uconfig.New[MapConfig](p)
	parsed, err := conf.Parse()
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	// 5. Verify parsed results
	expected := map[string]App{
		"auth.json":    {Name: "auth", Port: 8080},
		"gateway.json": {Name: "gateway", Port: 9000},
	}

	if !reflect.DeepEqual(parsed.Apps, expected) {
		t.Errorf("expected parsed apps to be %v, got %v", expected, parsed.Apps)
	}
}

func TestMapSliceConfig(t *testing.T) {
	// 1. Create a temporary directory
	tempDir, err := os.MkdirTemp("", "fileset_map_slice_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 2. Write test JSON files (slice per file for map[string][]App)
	app1Content := `[
		{"name": "auth", "port": 8080},
		{"name": "users", "port": 8081}
	]`
	app2Content := `[
		{"name": "gateway", "port": 9000}
	]`

	err = os.WriteFile(filepath.Join(tempDir, "auth.json"), []byte(app1Content), 0644)
	if err != nil {
		t.Fatalf("failed to write auth.json: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "gateway.json"), []byte(app2Content), 0644)
	if err != nil {
		t.Fatalf("failed to write gateway.json: %v", err)
	}

	// 3. Setup fileset plugin
	pattern := filepath.Join(tempDir, "*.json")
	p := New("my-apps", Absolute(pattern), json.Unmarshal)

	// 4. Load with uconfig
	conf := uconfig.New[MapSliceConfig](p)
	parsed, err := conf.Parse()
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	// 5. Verify parsed results
	expected := map[string][]App{
		"auth.json": {
			{Name: "auth", Port: 8080},
			{Name: "users", Port: 8081},
		},
		"gateway.json": {
			{Name: "gateway", Port: 9000},
		},
	}

	if !reflect.DeepEqual(parsed.Apps, expected) {
		t.Errorf("expected parsed apps to be %v, got %v", expected, parsed.Apps)
	}
}

