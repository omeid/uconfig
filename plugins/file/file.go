package file

import "io/ioutil"

// Unmarshal is any function that maps the source bytes to the provided
// config.
type Unmarshal func(src []byte, v interface{}) error

// File is a file loader plugin for uconfig
type File interface {
	Walk(conf interface{}) error
	Parse() error
}

// New returns an EnvSet.
func New(path string, unmarshal Unmarshal) File {
	return &walker{
		path:      path,
		unmarshal: unmarshal,
	}
}

type walker struct {
	path      string
	conf      interface{}
	unmarshal Unmarshal
}

func (v *walker) Walk(conf interface{}) error {
	v.conf = conf

	// check file?

	return nil
}

func (v *walker) Parse() error {

	src, err := ioutil.ReadFile(v.path)
	if err != nil {
		return err
	}

	return v.unmarshal(src, v.conf)
}
