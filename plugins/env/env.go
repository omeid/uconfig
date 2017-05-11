// Package env provides environment variables support for uconfig
package env

import (
	"os"
	"strings"

	"github.com/omeid/uconfig/flat"
)

const tag = "env"

// Envs is an env variable plugin.
type Envs interface {
	Visit(flat.Fields) error

	Parse() error
}

// New returns an EnvSet.
func New() Envs {
	return &visitor{}
}

type visitor struct {
	fields flat.Fields
}

func makeEnvName(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.ToUpper(name)

	return name
}
func (v *visitor) Visit(f flat.Fields) error {

	v.fields = f

	for _, f := range v.fields {
		name, ok := f.Tag(tag)
		if name == "-" {
			continue
		}

		if !ok || name == "" {
			name = makeEnvName(f.Name())
		}

		f.Meta()[tag] = "$" + name
	}

	return nil
}

func (v *visitor) Parse() error {

	for _, f := range v.fields {
		// Next block could use field.Meta and grab the tag name.
		/* start */
		name, ok := f.Tag(tag)
		if name == "-" {
			continue
		}

		if !ok || name == "" {
			name = makeEnvName(f.Name())
		}

		/* end */
		value := os.Getenv(name)

		if value == "" {
			continue
		}
		err := f.Set(value)
		if err != nil {
			return err
		}
	}

	return nil
}
