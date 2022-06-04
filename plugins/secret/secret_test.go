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
	Password   string `secret:""`
	EmptyValue string `secret:""`
	Alt        string `secret:"AltPassword"`
	Ignore     string

	Nested fNested
}

func TestSecret(t *testing.T) {

	expect := fSecrets{
		Password:   "password",
		Alt:        "altPassword",
		EmptyValue: "",
		Nested: fNested{
			Pass: "sub-pass",
		},
	}

	secrets := map[string]string{
		"PASSWORD":    "password",
		"AltPassword": "altPassword",
		"NESTED_PASS": "sub-pass",
		"EMPTYVALUE":  "",
	}

	source := func(name string) (string, error) {

		secret, ok := secrets[name]
		if !ok {
			return "", fmt.Errorf("couldn't find secret for %s", name)
		}
		return secret, nil
	}

	value := fSecrets{}
	conf, err := uconfig.New(&value, secret.New(source))
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

func TestSecretErr(t *testing.T) {

	source := func(name string) (string, error) {

		return "", fmt.Errorf("No value found for: %s", name)
	}

	value := fSecrets{}
	conf, err := uconfig.New(&value, secret.New(source))

	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()

	if err == nil {
		t.Fatal("Expected error but got nil")
		return
	}

	expect := "No value found for: PASSWORD"
	if err.Error() != expect {
		t.Fatalf("Expected: %s\nGot: %s", expect, err)
	}

}

func TestSecretSetErr(t *testing.T) {

	source := func(name string) (string, error) {
		return "not a number", nil
	}

	value := struct {
		Count int `secret:""`
	}{}

	conf, err := uconfig.New(&value, secret.New(source))

	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()

	if err == nil {
		t.Fatal("Expected error but got nil")
		return
	}

	expect := "strconv.ParseInt: parsing \"not a number\": invalid syntax"
	if err.Error() != expect {
		t.Fatalf("Expected: %s\nGot: %s", expect, err)
	}

}
