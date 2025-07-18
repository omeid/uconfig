package flag_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins/defaults"
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

	expect := &f.Config{
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

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[f.Config](fs)

	value, err := conf.Parse()
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

	expect := &fFlag{
		Address: "https://blah.bleh",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fFlag](fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

type fCli struct {
	Command string `flag:",command"`
	Address string `flag:"host"`
}

func TestFlagTagCommand(t *testing.T) {
	args := []string{
		"-host=https://blah.bleh",
		"run",
	}

	expect := &fCli{
		Command: "run",
		Address: "https://blah.bleh",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCli](fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

type fCliRename struct {
	Mode    string `flag:",command"`
	Command string `flag:"command"`
	Address string `flag:"host"`
}

func TestFlagTagCommandRename(t *testing.T) {
	args := []string{
		"-command=dance",
		"disco",
	}

	expect := &fCliRename{
		Mode:    "disco",
		Command: "dance",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliRename](fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

type fCliDefault struct {
	Mode string `flag:",command" default:"run"`
}

func TestFlagTagCommandDefault(t *testing.T) {
	args := []string{}

	expect := &fCliDefault{
		Mode: "run",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliDefault](defaults.New(), fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestFlagTagCommandDefaultOverride(t *testing.T) {
	args := []string{"walk"}

	expect := &fCliDefault{
		Mode: "walk",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliDefault](defaults.New(), fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestFlagTagCommandDefaultOverrideEmpty(t *testing.T) {
	args := []string{""}

	expect := &fCliDefault{
		Mode: "run",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliDefault](defaults.New(), fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestFlagTagCommandMissing(t *testing.T) {
	args := []string{""}

	expect := &fCliDefault{
		Mode: "run",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliDefault](fs, defaults.New())

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

type fCliRequired struct {
	Command string `flag:",required,command"`
	Mode    string `flag:",required"`
}

func TestFlagRequiredMissing(t *testing.T) {
	args := []string{}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliRequired](fs)
	_, err := conf.Parse()

	if err == nil {
		t.Fatal("expected error for missing required failed got nil")
	}

	expect := "missing required flag: [command]\nmissing required flag: mode"

	if err.Error() != expect {
		t.Errorf("expected (%s) but got (%s)", expect, err)
	}
}

func TestFlagRequiredMissingCommand(t *testing.T) {
	args := []string{"-mode", "slow"}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliRequired](fs)
	_, err := conf.Parse()

	if err == nil {
		t.Fatal("expected error for missing required failed got nil")
	}

	expect := "missing required flag: [command]"

	if err.Error() != expect {
		t.Errorf("expected (%s) but got (%s)", expect, err)
	}
}

type fCliBool struct {
	Command string `flag:",command"`
	Fast    bool
}

func TestFlagCommandAfterBool(t *testing.T) {
	args := []string{"-fast", "jump"}
	expect := &fCliBool{
		Fast:    true,
		Command: "jump",
	}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliBool](fs)

	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestFlagRequiredOkay(t *testing.T) {
	args := []string{"-mode=happy", "run"}

	expect := &fCliRequired{Mode: "happy", Command: "run"}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliRequired](fs)
	value, err := conf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}

func TestFlagExtraArgs(t *testing.T) {
	args := []string{"-mode=happy", "run", "fun"}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf := uconfig.New[fCliRequired](fs)

	_, err := conf.Parse()

	if err == nil {
		t.Fatal("Expected error for extra arguments but got nil")
	}

	expect := "extra arguments provided: (run)"

	if err.Error() != expect {
		t.Errorf("expected (%s) but got (%s)", expect, err)
	}
}
