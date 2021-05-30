// Package env provides environment variables support for uconfig
package env

import (
	"os"
	"strings"

	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/plugins"
)

const tag = "env"

func init() {
	plugins.RegisterTag(tag)
}

// New returns an EnvSet.
func New() plugins.Plugin {
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

		if !ok || name == "" {
			name = makeEnvName(f.Name())
		}

		f.Meta()[tag] = name
	}

	return nil
}

func (v *visitor) Parse() error {

	for _, f := range v.fields {
		name, ok := f.Meta()[tag]
		if !ok || name == "-" {
			continue
		}

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
