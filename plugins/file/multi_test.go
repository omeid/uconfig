package file_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
)

func TestMultiFile(t *testing.T) {
	expect := &f.Config{
		Command: "",
		Anon:    f.Anon{},

		GoHard: false,

		Redis: f.Redis{
			Host: "",
			Port: 0,
		},

		Rethink: f.RethinkConfig{
			Host: f.Host{
				Address: "rethink-cluster",
				Port:    "28015",
			},
			Db: "base",
		},
	}

	unmarshalOptions := file.UnmarshalOptions{
		".json": json.Unmarshal,
	}

	conf := uconfig.New[f.Config](
		file.NewMulti("testdata/config_rethink.json", unmarshalOptions, false),
	)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestMultiFileChoice(t *testing.T) {
	expect := &f.Config{
		Command: "",
		Anon:    f.Anon{},

		GoHard: false,

		Redis: f.Redis{
			Host: "",
			Port: 0,
		},

		Rethink: f.RethinkConfig{
			Host: f.Host{
				Address: "rethink-cluster",
				Port:    "28015",
			},
			Db: "base",
		},
	}

	unmarshalOptions := file.UnmarshalOptions{
		".json": json.Unmarshal,
		".conf": json.Unmarshal,
	}

	options := []string{
		"testdata/config_rethink.json",
		"testdata/rethinkdb/.conf",
	}
	for _, filepath := range options {

		conf := uconfig.New[f.Config](
			file.NewMulti(filepath, unmarshalOptions, false),
		)

		value, err := conf.Parse()
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expect, value); diff != "" {
			t.Error(diff)
		}
	}
}

func TestMultiOptional(t *testing.T) {
	expect := &f.Config{
		Command: "",
		Anon:    f.Anon{},

		GoHard: false,

		Redis: f.Redis{
			Host: "",
			Port: 0,
		},

		Rethink: f.RethinkConfig{
			Host: f.Host{
				Address: "rethink-cluster",
				Port:    "28015",
			},
			Db: "base",
		},
	}

	unmarshalOptions := file.UnmarshalOptions{
		".json": json.Unmarshal,
	}

	conf := uconfig.New[f.Config](
		file.NewMulti("testdata/doesnt_exists.json", unmarshalOptions, true),
		file.NewMulti("testdata/config_rethink.json", unmarshalOptions, false),
	)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}
