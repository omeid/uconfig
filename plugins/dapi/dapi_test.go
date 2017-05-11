package dapi_test

import (
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/dapi"
)

type fDAPIVersion struct {
	Build string `dapi:"annotations:build"`
}
type fDAPI struct {
	Version fDAPIVersion
	Builder string `dapi:"annotations:builder"`

	Cluster    string `dapi:"labels:cluster"`
	RackNumber string `dapi:"labels:rack"`
	AZ         string `dapi:"labels:zone"`
}

func TestEnvTag(t *testing.T) {

	envs := map[string]string{
		"MY_HOST_NAME": "https://blah.bleh",
	}

	for key, value := range envs {
		os.Setenv(key, value)
	}

	expect := fDAPI{
		Version: fDAPIVersion{
			Build: "two",
		},
		Builder: "john-doe",

		Cluster:    "test-cluster1",
		RackNumber: "rack-22",
		AZ:         "us-est-coast",
	}

	value := fDAPI{}

	conf, err := uconfig.New(&value)
	if err != nil {
		t.Fatal(err)
	}

	err = conf.Visitor(dapi.New("testdata"))
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
