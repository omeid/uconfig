package uconfig_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/omeid/uconfig/plugins/secret"
)

func TestLoadBasic(t *testing.T) {

	expect := f.Config{
		Command: "run",
		Anon: f.Anon{
			Version: "version-from-env",
		},

		GoHard: true,

		Redis: f.Redis{
			Host: "from-envs",
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

	files := uconfig.Files{
		{"testdata/classic.json", json.Unmarshal, true},
	}

	value := f.Config{}

	// set some env vars to test env var and plugin orders.
	os.Setenv("VERSION", "version-from-env")
	os.Setenv("REDIS_ADDRESS", "from-envs")
	// patch the os.Args. for our tests.

	_, err := uconfig.Load(&value, files)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}

func TestLoadWithSecret(t *testing.T) {

	//Config is part of text fixtures.
	type Creds struct {
		APIKey   string `secret:""`
		APIToken string `secret:"API_TOKEN"`
	}

	type Config struct {
		Redis   f.Redis
		Rethink f.RethinkConfig
		Creds   Creds
	}
	expect := Config{
		Redis: f.Redis{
			Host: "redis-host",
			Port: 6379,
		},

		Rethink: f.RethinkConfig{
			Host: f.Host{
				Address: "rethink-cluster",
				Port:    "28015",
			},
			Db:       "base",
			Password: "top secret token",
		},

		Creds: Creds{
			APIKey:   "top secret token",
			APIToken: "top secret token",
		},
	}

	files := uconfig.Files{
		{"testdata/classic.json", json.Unmarshal, true},
	}

	value := Config{}

	SecretProvider := func(name string) (string, error) {

		// known secrets.
		if name == "API_TOKEN" || name == "RETHINK_PASSWORD" || name == "CREDS_APIKEY" {
			return "top secret token", nil
		}

		return "", fmt.Errorf("Secret not found %s", name)
	}

	os.Unsetenv("VERSION")
	os.Unsetenv("REDIS_ADDRESS")

	_, err := uconfig.Load(&value, files, secret.New(SecretProvider))
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}

func TestLoadWithMultiFile(t *testing.T) {

	expect := f.Config{
		Command: "run",
		Anon: f.Anon{
			Version: "version-from-env",
		},

		Redis: f.Redis{
			Host: "from-envs",
		},

		Rethink: f.RethinkConfig{
			Host: f.Host{
				Address: "rethink-cluster",
				Port:    "28015",
			},
			Db: "base",
		},
	}

	options := uconfig.UnmarshalOptions{
		".json": json.Unmarshal,
	}

	value := f.Config{}

	// set some env vars to test env var and plugin orders.
	os.Setenv("VERSION", "version-from-env")
	os.Setenv("REDIS_ADDRESS", "from-envs")
	// patch the os.Args. for our tests.

	_, err := uconfig.Load(&value,
		nil,
		file.NewMulti("plugins/file/testdata/config_rethink.json", options, true),
	)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}
func TestLoadBadPlugin(t *testing.T) {

	var badPlugin BadPlugin

	config := f.Config{}

	_, err := uconfig.Load(&config, nil, badPlugin)

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "Unsupported plugins. Expecting a Walker or Visitor" {
		t.Errorf("Expected unsupported plugin error, got: %v", err)
	}

}
