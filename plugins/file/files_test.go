package file_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
)

func TestFiles(t *testing.T) {
	expect := &f.Config{
		Command: "run-file",
		Anon: f.Anon{
			Version: "0.2",
		},

		GoHard: true,

		Redis: f.Redis{
			Host: "redis-host",
			Port: 6379,
		},

		Rethink: f.RethinkConfig{
			Host: f.Host{
				Address: "rethink-cluster",
				Port:    "28015",
			},
			Db: "base",
		},
	}

	files := file.Files{
		{"testdata/config_rethink.json", json.Unmarshal, true},
		{"testdata/config_partial.json", json.Unmarshal, true},
		{"testdata/.local.json", json.Unmarshal, false},
	}

	os.Args = os.Args[:1]
	conf := uconfig.Classic[f.Config](files)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}
