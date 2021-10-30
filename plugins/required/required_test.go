package required_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/required"
)

type fRequired struct {
	Required    string `required:"true"`
	NotRequired string
	NotComparable fmt.Stringer `required:"true"`
}

func TestRequired(t *testing.T) {
	expect := fRequired{NotRequired: "not-empty"}

	conf, err := uconfig.New(&expect, required.New())
	if err != nil {
		t.Fatal(err)
	}

	if err = conf.Parse(); !errors.As(err, &required.ErrRequiredField{}) {
		t.Fatalf("expected Parse() to fail with required field error but did not, instead received: %#v", err)
	}

	expect.Required = "something"
	if err := conf.Parse(); err != nil {
		t.Fatalf("expected Parse() to succeed but received: %#v", err)
	}
}
