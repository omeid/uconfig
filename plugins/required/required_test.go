package required_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/required"
)

type fRequired struct {
	Required      string `required:"true"`
	NotRequired   string
	NotComparable fmt.Stringer `required:"true"`
}

func TestRequired(t *testing.T) {
	expect := fRequired{NotRequired: "not-empty"}

	conf, err := uconfig.New(&expect, required.New())
	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()
	t.Log(err)

	var fieldError *required.ErrRequiredField

	if !errors.As(err, &fieldError) {
		t.Fatalf("expected Parse() to fail with required field error but did not, instead received: %#v", err)
	}

	fieldErrorName := fieldError.Name()
	if fieldErrorName != "Required" {
		t.Fatalf("Expected to fail on first field `Required` but failed on `%s`", fieldErrorName)
	}

	expect.Required = "something"
	if err := conf.Parse(); err != nil {
		t.Fatalf("expected Parse() to succeed but received: %#v", err)
	}
}
