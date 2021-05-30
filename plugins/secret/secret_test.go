package secret_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/secret"
)

type fNested struct {
	Pass   string `secret:""`
	Ignore string
}

type fSecrets struct {
	Password string `secret:""`
	Alt      string `secret:"AltPassword"`
	Ignore   string

	Nested fNested
}

func TestDefaultTag(t *testing.T) {

	expect := fSecrets{
		Password: "password",
		Alt:      "altPassword",
		Nested: fNested{
			Pass: "sub-pass",
		},
	}

	value := fSecrets{}

	conf, err := uconfig.New(&value)
	if err != nil {
		t.Fatal(err)
	}

	secrets := map[string]string{
		"PASSWORD":    "password",
		"AltPassword": "altPassword",
		"NESTED_PASS": "sub-pass",
	}

	source := func(name string) (string, error) {

		secret, ok := secrets[name]
		if !ok {
			return "", fmt.Errorf("couldn't find secret for %s", name)
		}
		return secret, nil
	}

	err = conf.AddPlugin(secret.New(source))
	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}
