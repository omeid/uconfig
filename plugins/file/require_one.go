package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/omeid/uconfig/plugins"
)

var (
	ErrNoFiles       = errors.New("require one: no files found")
	ErrMultipleFiles = errors.New("require one: multiple files found")
)

// RequireOne represents a set of file paths and their unmarshal functions,
// where exactly one of the files must exist.
type RequireOne []struct {
	Path      Source
	Unmarshal Unmarshal
}

// Plugins constructs a slice of Plugin from the RequireOne list.
func (r RequireOne) Plugins() []plugins.Plugin {
	return []plugins.Plugin{&requireOneWalker{files: r}}
}

type requireOneWalker struct {
	files RequireOne
	match *struct {
		Path      Source
		Unmarshal Unmarshal
	}
	matchResolved string
	conf          any
}

func (e *requireOneWalker) Walk(conf any) error {
	e.conf = conf
	var found []string

	for i, f := range e.files {
		path := f.Path.Resolve()

		_, err := os.Stat(path)
		if err == nil {
			found = append(found, path)
			if e.match == nil {
				e.match = &e.files[i]
				e.matchResolved = path
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	if len(found) == 0 {
		return ErrNoFiles
	}
	if len(found) > 1 {
		return fmt.Errorf("%w: %v", ErrMultipleFiles, found)
	}

	return nil
}

func (e *requireOneWalker) Parse() error {
	if e.match == nil {
		return ErrNoFiles
	}

	f, err := os.Open(e.matchResolved)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = e.match.Unmarshal(data, e.conf)
	if err != nil {
		filePath := errors.New(e.matchResolved)
		return errors.Join(filePath, err)
	}

	return nil
}

// SourcePaths implements plugins.SourcePaths.
func (e *requireOneWalker) SourcePaths() []plugins.SourcePath {
	if e.match != nil {
		return []plugins.SourcePath{{Name: e.match.Path.Name(), Path: e.matchResolved}}
	}
	return nil
}

// Usage implements plugins.Usage.
func (e *requireOneWalker) Usage(fieldname string) (string, string) {
	if fieldname == "." {
		var pathsList []string
		for _, f := range e.files {
			if f.Path.Kind() != KindUnknown {
				pathsList = append(pathsList, fmt.Sprintf("%-10s %s", f.Path.Kind().String()+":", f.Path.Name()))
			} else {
				pathsList = append(pathsList, f.Path.Name())
			}
		}
		return "Files Required One", strings.Join(pathsList, "\n")
	}
	return "", ""
}
