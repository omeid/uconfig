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
	expect := &f.Config{
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

	// set some env vars to test env var and plugin orders.
	err := os.Setenv("VERSION", "version-from-env")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("REDIS_ADDRESS", "from-envs")
	if err != nil {
		t.Fatal(err)
	}
	// patch the os.Args. for our tests.

	conf := uconfig.Load[f.Config](files)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestLoadWithSecret(t *testing.T) {
	// Config is part of text fixtures.
	type Creds struct {
		APIKey   string `secret:""`
		APIToken string `secret:"API_TOKEN"`
	}

	type Config struct {
		Redis   f.Redis
		Rethink f.RethinkConfig
		Creds   Creds
	}
	expect := &Config{
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

	SecretProvider := func(name string) (string, error) {
		// known secrets.
		if name == "API_TOKEN" || name == "RETHINK_PASSWORD" || name == "CREDS_APIKEY" {
			return "top secret token", nil
		}

		return "", fmt.Errorf("Secret not found %s", name)
	}

	err := os.Unsetenv("VERSION")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Unsetenv("REDIS_ADDRESS")
	if err != nil {
		t.Fatal(err)
	}

	conf := uconfig.Load[Config](files, secret.New(SecretProvider))

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestLoadWithMultiFile(t *testing.T) {
	expect := &f.Config{
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

	// set some env vars to test env var and plugin orders.
	err := os.Setenv("VERSION", "version-from-env")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("REDIS_ADDRESS", "from-envs")
	if err != nil {
		t.Fatal(err)
	}

	conf := uconfig.Load[f.Config](
		nil,
		file.NewMulti("plugins/file/testdata/config_rethink.json", options, true),
	)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestLoadBadPlugin(t *testing.T) {
	var badPlugin BadPlugin

	conf := uconfig.Load[f.Config](nil, badPlugin)
	_, err := conf.Parse()

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "unsupported plugins. expecting a walker or visitor" {
		t.Errorf("Expected unsupported plugin error, got: %v", err)
	}
}
