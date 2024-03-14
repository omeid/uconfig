// Package secret enable uconfig to integrate with secret plugins.
package secret

import (
	"fmt"
	"strings"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

const tag = "secret"

func init() {
	plugins.RegisterTag(tag)
}

// Sourcer is any function that can exchanges a secret name with it's value.
type Sourcer func(string) (string, error)

// New returns the secret provider.
func New(source Sourcer) plugins.Plugin {
	return &secret{source: source}
}

type secret struct {
	fields flat.Fields
	source Sourcer
}

func makeSecretName(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.ToUpper(name)

	return name
}

func (v *secret) Visit(f flat.Fields) error {

	v.fields = f

	for _, f := range v.fields {

		// secret only works with explicitly tagged.
		if _, ok := f.Tag(tag); !ok {
			continue
		}

		name, explicit := f.Name(tag)
		if !explicit {
			name = makeSecretName(name)
		}

		fmt.Printf("Secret name %s %t\n", name, explicit)
		f.Meta()[tag] = name
	}

	return nil
}

func (v *secret) Parse() error {

	for _, f := range v.fields {
		name := f.Meta()[tag]
		if name == "" || name == "-" {
			continue
		}

		n, _ := f.Name(tag)
		fmt.Printf("Looking up password for %s: %s\n", name, n)
		value, err := v.source(name)

		if err != nil {
			return err
		}

		if value == "" {
			continue
		}

		err = f.Set(value)
		if err != nil {
			return err
		}
	}

	return nil
}
