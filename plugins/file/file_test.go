package file_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/go-test/deep"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
)

func TestEnvBasic(t *testing.T) {

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

	srcTOML := `
	Version = "0.2"

	GoHard =  true

	[Redis]
	Host =  "redis-host"
	Port =  6379

	[Rethink]
		Db =  "base"

		[Rethink.Host]
			Address =  "rethink-cluster"
			Port =     "28015"
	`

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
			"toml",
			bytes.NewReader([]byte(srcTOML)),
			toml.Unmarshal,
		},
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

		err = conf.Walker(file.NewReader(tc.Source, tc.Unmarshall))
		if err != nil {
			t.Fatal(err)
		}

		err = conf.Parse()

		if err != nil {
			t.Fatal(err)
		}

		if diff := deep.Equal(expect, value); diff != nil {
			t.Error(diff)
		}

	}
}
