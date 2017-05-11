package defaults_test

import (
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/defaults"
)

type fDefaults struct {
	Address string `default:"https://blah.bleh"`
}

func TestDefaultTag(t *testing.T) {

	envs := map[string]string{
		"MY_HOST_NAME": "https://blah.bleh",
	}

	for key, value := range envs {
		os.Setenv(key, value)
	}

	expect := fDefaults{
		Address: "https://blah.bleh",
	}

	value := fDefaults{}

	conf, err := uconfig.New(&value)
	if err != nil {
		t.Fatal(err)
	}

	err = conf.Visitor(defaults.New())
	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()

	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(expect, value); diff != nil {
		t.Error(diff)
	}

}
