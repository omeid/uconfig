package uconfig

import "github.com/omeid/uconfig/plugins/defaults"

// Defaults is a simple helper that setups a uconfig with default plugin and parses
// the defaults. Useful in tests.
func Defaults(conf interface{}) error {
	confs, err := New(&conf)
	if err != nil {
		return err
	}

	err = confs.Visitor(defaults.New())
	if err != nil {
		return err
	}

	return confs.Parse()
}
