package env_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/env"
)

func TestEnvBasic(t *testing.T) {
	envs := map[string]string{
		"COMMAND":              "start",
		"GOHARD":               "T",
		"VERSION":              "0.2",
		"REDIS_ADDRESS":        "redis-host",
		"REDIS_PORT":           "6379",
		"RETHINK_HOST_ADDRESS": "rethink-cluster",
		"RETHINK_HOST_PORT":    "28015",
		"RETHINK_DB":           "",
	}

	expect := &f.Config{
		Command: "start",
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
			Db: "",
		},
	}

	for key, value := range envs {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatal(err)
		}
	}

	conf := uconfig.New[f.Config](env.New())

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

type fEnv struct {
	Address string `env:"MY_HOST_NAME"`
}

func TestEnvTag(t *testing.T) {
	envs := map[string]string{
		"MY_HOST_NAME": "https://blah.bleh",
	}

	for key, value := range envs {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatal(err)
		}
	}

	expect := &fEnv{
		Address: "https://blah.bleh",
	}

	conf := uconfig.New[fEnv](env.New())

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}
