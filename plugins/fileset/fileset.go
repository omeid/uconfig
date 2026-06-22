package fileset

import (
	"fmt"
	"io"
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

type visitor struct {
	name         string
	set          Set
	unmarshal    file.Unmarshal
	currentData  []byte
	fields       flat.Fields
	matchedField flat.Field
}

// New returns a fileset plugin.
func New(name string, set Set, unmarshal file.Unmarshal) plugins.Plugin {
	v := &visitor{
		name:      name,
		set:       set,
		unmarshal: unmarshal,
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

	sources, err := v.set.Resolve()
	if err != nil {
		return fmt.Errorf("fileset: failed to resolve path %q: %w", v.set.Name(), err)
	}

	for _, match := range sources {
		rc, err := match.Open()
		if err != nil {
			return fmt.Errorf("fileset: failed to open file %q: %w", match.Resolve(), err)
		}

		data, err := io.ReadAll(rc)
		closeErr := rc.Close()
		if err != nil {
			return fmt.Errorf("fileset: failed to read file %q: %w", match.Resolve(), err)
		}
		if closeErr != nil {
			return fmt.Errorf("fileset: failed to close file %q: %w", match.Resolve(), closeErr)
		}

		v.currentData = data
		key := filepath.Base(match.Resolve())
		err = v.matchedField.Append(key, func(ptr any) error {
			return v.unmarshal(v.currentData, ptr)
		})
		if err != nil {
			return fmt.Errorf("fileset: failed to append file %q: %w", match.Resolve(), err)
		}
	}

	return nil
}

// SourcePaths implements plugins.SourcePaths.
func (v *visitor) SourcePaths() []plugins.SourcePath {
	pattern := v.set.Name()
	if pattern == "" {
		return nil
	}

	if idx := strings.Index(pattern, ": "); idx != -1 {
		pattern = pattern[idx+2:]
	}

	baseDir := getBaseDir(pattern)
	absDir, err := filepath.Abs(baseDir)
	if err == nil {
		return []plugins.SourcePath{{Name: v.set.Name(), Path: absDir}}
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
		kind := v.set.Kind()
		if kind != file.KindUnknown {
			return "Fileset", kind.String() + ": " + v.set.Name() + "\t\n"
		}
		return "Fileset", v.set.Name() + "\t\n"
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
