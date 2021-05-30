package file_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
)

func TestFileReader(t *testing.T) {

	expect := f.Config{
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

	srcJSON := `{
		"Version": "0.2",
		"GoHard": true,
		"Redis": {
			"Host": "redis-host",
			"Port": 6379
		},
		"Rethink": {
			"Db": "base",
			"Host": {
				"Address": "rethink-cluster",
				"Port": "28015"
			}
		}
	}`

	type TestCase struct {
		Name       string
		Source     io.Reader
		Unmarshall func([]byte, interface{}) error
	}

	for _, tc := range []TestCase{
		{
			"json",
			bytes.NewReader([]byte(srcJSON)),
			json.Unmarshal,
		},
	} {

		value := f.Config{}

		conf, err := uconfig.New(&value)
		if err != nil {
			t.Fatal(err)
		}

		err = conf.AddPlugin(file.NewReader(tc.Source, tc.Unmarshall))
		if err != nil {
			t.Fatal(err)
		}

		err = conf.Parse()

		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expect, value); diff != "" {
			t.Error(diff)
		}

	}
}

func TestFileOpen(t *testing.T) {

	expect := f.Config{
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

	type TestCase struct {
		Name       string
		Source     string
		Unmarshall func([]byte, interface{}) error
	}

	for _, tc := range []TestCase{
		{
			"json",
			"testdata/config.json",
			json.Unmarshal,
		},
	} {

		value := f.Config{}

		conf, err := uconfig.New(&value)
		if err != nil {
			t.Fatal(err)
		}

		err = conf.AddPlugin(file.New(tc.Source, tc.Unmarshall))
		if err != nil {
			t.Fatal(err)
		}

		err = conf.Parse()

		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expect, value); diff != "" {
			t.Error(diff)
		}

	}
}

func TestMulti(t *testing.T) {

	expect := f.Config{
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

	srcJSON := `{
		"Version": "0.2",
		"GoHard": true,
		"Redis": {
			"Host": "redis-host",
			"Port": 6379
		}
	}`

	reader := file.NewReader(bytes.NewReader([]byte(srcJSON)), json.Unmarshal)
	open := file.New("testdata/config_rethink.json", json.Unmarshal)

	value := f.Config{}
	conf, err := uconfig.New(&value, reader, open)
	if err != nil {
		t.Fatal(err)
	}
	err = conf.Parse()

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}

func TestBadFile(t *testing.T) {

	open, err := os.Open("testdata/config_rethink.json")
	if err != nil {
		t.Fatal(err)
	}

	open.Close() // close it so we get an error!
	reader := file.NewReader(open, json.Unmarshal)

	value := f.Config{}
	conf, err := uconfig.New(&value, reader)
	if err != nil {
		t.Fatal(err)
	}
	err = conf.Parse()

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if err.Error() != "read testdata/config_rethink.json: file already closed" {
		t.Errorf("Unexpected error: %v", err)
	}
}
