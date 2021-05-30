package plugins_test

import (
	"testing"

	"github.com/omeid/uconfig/plugins"
)

func TestRegisterTagMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected RegisterTag to panic on double register. %#v", r)
		}
	}()

	plugins.RegisterTag("duplicate")

	// registering the same tag twice should result into panic
	plugins.RegisterTag("duplicate")
}

func TestRegisterTagMustNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected RegisterTag not to panic: %v", r)
		}
	}()

	plugins.RegisterTag("uniuq.1")
	plugins.RegisterTag("uniuq.2")

}
