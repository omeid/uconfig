package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type absolute string

func (a absolute) Name() string    { return string(a) }
func (a absolute) Kind() Kind      { return KindAbsolute }
func (a absolute) Resolve() string { return string(a) }
func (a absolute) Open() (io.ReadCloser, error) {
	return os.Open(string(a))
}

// Absolute returns a Source for a fixed absolute path.
func Absolute(path string) Source {
	if path == "-" {
		return Stdin()
	}
	return absolute(path)
}

type relative string

func (r relative) Name() string { return string(r) }
func (r relative) Kind() Kind   { return KindRelative }
func (r relative) Resolve() string {
	abs, err := filepath.Abs(string(r))
	if err != nil {
		return string(r)
	}
	return abs
}

func (r relative) Open() (io.ReadCloser, error) {
	return os.Open(r.Resolve())
}

// Relative returns a Source that resolves a relative path against
// the working directory at the time of the call.
func Relative(path string) Source {
	if path == "-" {
		return Stdin()
	}
	return relative(path)
}

type workspace string

func (w workspace) Name() string { return string(w) }
func (w workspace) Kind() Kind   { return KindWorkspace }
func (w workspace) Resolve() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		candidate := filepath.Join(dir, string(w))
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return ""
		}
		dir = parent
	}
}

func (w workspace) Open() (io.ReadCloser, error) {
	return os.Open(w.Resolve())
}

// Workspace returns a Source that walks up the directory tree looking
// for a file at the given relative path (e.g. ".myapp/config").
func Workspace(name string) Source {
	return workspace(name)
}

type stdin struct {
	once sync.Once
	data []byte
	err  error
}

func (w *stdin) Name() string    { return "-" }
func (s *stdin) Kind() Kind      { return KindStdin }
func (w *stdin) Resolve() string { return "-" }
func (w *stdin) Open() (io.ReadCloser, error) {
	w.once.Do(func() {
		w.data, w.err = io.ReadAll(os.Stdin)
	})
	if w.err != nil {
		return nil, w.err
	}
	return io.NopCloser(bytes.NewReader(w.data)), nil
}

// Stdin returns a Source that reads from Stdin.
func Stdin() Source {
	return &stdin{}
}

type readerPath struct {
	name string
	r    io.Reader
}

func (r readerPath) Name() string    { return r.name }
func (r readerPath) Kind() Kind      { return KindUnknown }
func (r readerPath) Resolve() string { return r.name }
func (r readerPath) Open() (io.ReadCloser, error) {
	if rc, ok := r.r.(io.ReadCloser); ok {
		return rc, nil
	}
	return io.NopCloser(r.r), nil
}

// Reader returns a Source that reads from the provided io.Reader.
func Reader(src io.Reader, name string) Source {
	return readerPath{name: name, r: src}
}
