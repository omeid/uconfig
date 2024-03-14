package flag_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/flag"
)

func TestFlagBasic(t *testing.T) {

	args := []string{
		"-gohard",
		"-version=0.2",
		"-redis-address=redis-host",
		"-redis-port=6379",
		"-rethink-host-address=rethink-cluster",
		"-rethink-host-port=28015",
		"-rethink-db=base",
	}

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

	value := f.Config{}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf, err := uconfig.New(&value, fs)
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

type fFlag struct {
	Address string `flag:"host"`
}

func TestFlagTag(t *testing.T) {

	args := []string{
		"-host=https://blah.bleh",
	}

	expect := fFlag{
		Address: "https://blah.bleh",
	}

	value := fFlag{}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf, err := uconfig.New(&value, fs)
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
