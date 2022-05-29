package uconfig_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/secret"
)

func TestMain(m *testing.M) {
	// for go test framework.
	flag.Parse()

	os.Exit(m.Run())
}

func TestClassicBasic(t *testing.T) {

	expect := f.Config{
		Anon: f.Anon{
			Version: "from-flags",
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
	os.Setenv("VERSION", "bad-value-overrided-with-flags")
	os.Setenv("REDIS_HOST", "from-envs")
	// patch the os.Args. for our tests.
	os.Args = append(os.Args[:1], "-version=from-flags")

	_, err := uconfig.Classic(&value, files)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}

func TestClassicWithSecret(t *testing.T) {

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

	// patch the os.Args. for our tests.
	os.Args = os.Args[:1]
	os.Unsetenv("REDIS_HOST")

	_, err := uconfig.Classic(&value, files, secret.New(SecretProvider))
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}

}

func TestClassicBadPlugin(t *testing.T) {

	var badPlugin BadPlugin

	config := f.Config{}

	_, err := uconfig.Classic(&config, nil, badPlugin)

	if err == nil {
		t.Error("expected error for bad plugin, got nil")
	}

	if err.Error() != "Unsupported plugins. Expecting a Walker or Visitor" {
		t.Errorf("Expected unsupported plugin error, got: %v", err)
	}

}
