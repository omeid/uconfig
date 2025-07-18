package uconfig_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
)

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Must should not panic, but did: %v", r)
		}
	}()

	uconfig.Must[f.Config]()
}

func TestMustPanic(t *testing.T) {
	defer func() {
		r := recover()

		if r == nil {
			t.Error("Was expecting panic but got nil")
			return
		}

		expectErr := "read testdata/classic.json: file already closed"

		if err, ok := r.(error); !ok || err.Error() != expectErr {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	filepath := "testdata/classic.json"
	open, err := os.Open(filepath)
	if err != nil {
		t.Fatal(err)
	}

	err = open.Close() // close it so we get an error!
	if err != nil {
		t.Fatal(err)
	}
	badFile := file.NewReader(open, filepath, json.Unmarshal)

	var buf bytes.Buffer
	uconfig.UsageOutput = &buf

	uconfig.Must[f.Config](badFile)
}
