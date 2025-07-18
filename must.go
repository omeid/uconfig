package uconfig

import "github.com/omeid/uconfig/plugins"

// Must is like New but also calls Parse and panics instead
// of returning errors. This is useful in tests.
func Must[C any](plugins ...plugins.Plugin) *C {
	c := New[C](plugins...)

	conf, err := c.Parse()
	if err != nil {
		c.Usage()
		panic(err)
	}

	return conf
}
