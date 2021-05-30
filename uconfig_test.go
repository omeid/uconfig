package uconfig_test

import (
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins"
)

type BadPlugin interface {
	plugins.Plugin

	NotWalkerOrVisitor()
}

func TestBadPlug(t *testing.T) {

	var badPlugin BadPlugin

	config := f.Config{}

	_, err := uconfig.New(&config, badPlugin)

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "Unsupported plugins. Expecting a Walker or Visitor" {
		t.Errorf("Expected unsupported plugin error, got: %v", err)
	}
}
