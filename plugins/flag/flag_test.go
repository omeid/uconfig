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

type fCli struct {
	Command string `flag:",command"`
	Address string `flag:"host"`
}

func TestFlagTagCommand(t *testing.T) {

	args := []string{
		"-host=https://blah.bleh",
		"run",
	}

	expect := fCli{
		Command: "run",
		Address: "https://blah.bleh",
	}

	value := fCli{}

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

	expect := fCliRename{
		Mode:    "disco",
		Command: "dance",
	}

	value := fCliRename{}

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

type fCliDefault struct {
	Mode string `flag:",command" default:"run"`
}

func TestFlagTagCommandDefault(t *testing.T) {

	args := []string{}

	expect := fCliDefault{
		Mode: "run",
	}

	value := fCliDefault{
		Mode: "run",
	}

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
func TestFlagTagCommandDefaultOverride(t *testing.T) {

	args := []string{"walk"}

	expect := fCliDefault{
		Mode: "walk",
	}

	value := fCliDefault{
		Mode: "fast",
	}

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

func TestFlagTagCommandDefaultOverrideEmpty(t *testing.T) {

	args := []string{""}

	expect := fCliDefault{
		Mode: "walk",
	}

	value := fCliDefault{
		Mode: "walk",
	}

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

type fCliRequired struct {
	Command string `flag:",required,command"`
	Mode    string `flag:",required"`
}

func TestFlagRequiredMissing(t *testing.T) {

	args := []string{}

	value := fCliRequired{}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf, err := uconfig.New(&value, fs)

	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()

	if err == nil {
		t.Fatal("Expected error for missing required failed got nil")
	}

	expect := "Missing required flag: [command]\nMissing required flag: mode"

	if err.Error() != expect {
		t.Errorf("expected (%s) but got (%s)", expect, err)
	}

}

func TestFlagRequiredMissingCommand(t *testing.T) {

	args := []string{"-mode=slow"}

	value := fCliRequired{}

	fs := flag.New("testing", flag.PanicOnError, args)

	conf, err := uconfig.New(&value, fs)

	if err != nil {
		t.Fatal(err)
	}

	err = conf.Parse()

	if err == nil {
		t.Fatal("Expected error for missing required failed got nil")
	}

	expect := "Missing required flag: [command]"

	if err.Error() != expect {
		t.Errorf("expected (%s) but got (%s)", expect, err)
	}

}

func TestFlagRequiredOkay(t *testing.T) {

	args := []string{"-mode=happy", "run"}

	value := fCliRequired{}
	expect := fCliRequired{Mode: "happy", Command: "run"}

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
