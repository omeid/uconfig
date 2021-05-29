package defaults_test

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/defaults"
)

type fDefaults struct {
	Address string        `default:"https://blah.bleh"`
	Bases   []string      `default:"list,blah"`
	Timeout time.Duration `default:"5s"`
}

func TestDefaultTag(t *testing.T) {

	expect := fDefaults{
		Address: "https://blah.bleh",
		Bases:   []string{"list", "blah"},
		Timeout: 5 * time.Second,
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
