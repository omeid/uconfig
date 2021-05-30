package uconfig_test

import (
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
)

func TestUsage(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Usage should not panic, but did: %v", r)
		}
	}()

	value := f.Config{}
	c, err := uconfig.Classic(&value, nil)
	if err != nil {
		t.Fatal(err)
	}

	c.Usage()

}
