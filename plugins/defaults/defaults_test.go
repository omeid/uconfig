package defaults_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/defaults"
)

type fDefaults struct {
	Address string        `default:"https://blah.bleh"`
	Bases   []string      `default:"list,blah"`
	Timeout time.Duration `default:"5s"`
	Ignored string
}

func TestDefaultTag(t *testing.T) {
	expect := &fDefaults{
		Address: "https://blah.bleh",
		Bases:   []string{"list", "blah"},
		Timeout: 5 * time.Second,
		Ignored: "",
	}

	conf := uconfig.New[fDefaults](defaults.New())

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}
