// Package fileset provides a fileset plugin for uconfig.
package fileset

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/file"
)

const tag = "fileset"

func init() {
	plugins.RegisterTag(tag)
}

// Path pairs a display name with a lazy path resolver.
type Path struct {
	Name    string
	Kind    string // "absolute", "relative", etc.
	Resolve func() ([]string, error)
}

// Absolute returns a Path for a fixed absolute path/glob.
func Absolute(pattern string) Path {
	return Path{
		Name:    pattern,
		Kind:    "absolute",
		Resolve: func() ([]string, error) { return filepath.Glob(pattern) },
	}
}

// Relative returns a Path that resolves a relative path/glob against
// the working directory at the time of the call.
func Relative(pattern string) Path {
	return Path{
		Name: pattern,
		Kind: "relative",
		Resolve: func() ([]string, error) {
			abs, err := filepath.Abs(pattern)
			if err != nil {
				return nil, err
			}
			return filepath.Glob(abs)
		},
	}
}

type visitor struct {
	name         string
	path         Path
	unmarshal    func(any) error
	currentData  []byte
	fields       flat.Fields
	matchedField flat.Field
}

// New returns a fileset plugin.
func New(name string, path Path, unmarshaller file.Unmarshal) plugins.Plugin {
	v := &visitor{
		name: name,
		path: path,
	}
	v.unmarshal = func(ptr any) error {
		return unmarshaller(v.currentData, ptr)
	}
	return v
}

// VisitFolds implements plugins.FoldVisitor.
func (v *visitor) VisitFolds(fields flat.Fields) error {
	v.fields = fields
	var found flat.Field
	for _, f := range fields {
		tagValue, ok := f.Tag(tag)
		if ok && tagValue == v.name {
			if found != nil {
				return fmt.Errorf("fileset: multiple fields found with fileset tag %q", v.name)
			}
			found = f
		}
	}

	if found == nil {
		return fmt.Errorf("fileset: no field found with fileset tag %q", v.name)
	}

	v.matchedField = found
	return nil
}

// Parse implements plugins.Plugin.
func (v *visitor) Parse() error {
	if v.matchedField == nil {
		return fmt.Errorf("fileset: visit was not called or failed to find matching field for tag %q", v.name)
	}

	files, err := v.path.Resolve()
	if err != nil {
		return fmt.Errorf("fileset: failed to resolve path %q: %w", v.path.Name, err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("fileset: failed to read file %q: %w", file, err)
		}

		v.currentData = data
		key := filepath.Base(file)
		err = v.matchedField.Append(key, v.unmarshal)
		if err != nil {
			return fmt.Errorf("fileset: failed to append file %q: %w", file, err)
		}
	}

	return nil
}

// SourcePaths implements plugins.SourcePaths.
func (v *visitor) SourcePaths() []plugins.SourcePath {
	pattern := v.path.Name
	if pattern == "" {
		return nil
	}

	if idx := strings.Index(pattern, ": "); idx != -1 {
		pattern = pattern[idx+2:]
	}

	baseDir := getBaseDir(pattern)
	absDir, err := filepath.Abs(baseDir)
	if err == nil {
		return []plugins.SourcePath{{Name: v.path.Name, Path: absDir}}
	}
	return nil
}

// Usage implements plugins.Usage.
func (v *visitor) Usage(fieldname string) (string, string) {
	name := v.name
	if v.matchedField != nil {
		name, _ = v.matchedField.Name("")
	}
	if name == fieldname {
		return "Fileset", v.path.Kind + ": " + v.path.Name + "\t\n"
	}
	return "", ""
}

// getBaseDir returns the static directory prefix of a glob pattern.
func getBaseDir(pattern string) string {
	idx := -1
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		if c == '*' || c == '?' || c == '[' {
			idx = i
			break
		}
	}
	if idx == -1 {
		return filepath.Dir(pattern)
	}
	return filepath.Dir(pattern[:idx])
}
