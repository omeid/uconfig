package uconfig

import "github.com/omeid/uconfig/plugins"

// Must is like New but also calls Parse and panics instead
// of returning errors. This is useful in tests.
func Must(conf interface{}, plugins ...plugins.Plugin) {

	c, err := New(conf, plugins...)
	if err != nil {
		panic(err)
	}

	err = c.Parse()
	if err != nil {
		c.Usage()
		panic(err)
	}

}
