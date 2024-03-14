package uconfig_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/flat"
	"github.com/omeid/uconfig/internal/f"
	"github.com/omeid/uconfig/plugins"
	"github.com/omeid/uconfig/plugins/secret"
)

const expectedUsageMessage = `
Supported Fields:
FIELD                   FLAG                     ENV                     DEFAULT    GOODPLUGIN              SECRET              USAGE
-----                   -----                    -----                   -------    ----------              ------              -----
Version                 -version                 VERSION                            Version                                     
GoHard                  -gohard                  GOHARD                             GoHard                                      
Redis.Address           -redis-address           REDIS_ADDRESS                      Redis.Address                               
Redis.Port              -redis-port              REDIS_PORT                         Redis.Port                                  
Rethink.Host.Address    -rethink-host-address    RETHINK_HOST_ADDRESS               Rethink.Host.Address                        
Rethink.Host.Port       -rethink-host-port       RETHINK_HOST_PORT                  Rethink.Host.Port                           
Rethink.Db              -rethink-db              RETHINK_DB              primary    Rethink.Db                                  main database used by our application
Rethink.Password        -rethink-password        RETHINK_PASSWORD                   Rethink.Password        RETHINK_PASSWORD    
`

type UselessPluginVisitor struct {
	plugins.Plugin
}

func (*UselessPluginVisitor) Parse() error { return nil }

func (*UselessPluginVisitor) Visit(fields flat.Fields) error {
	for _, f := range fields {
		name, _ := f.Name("goodplugin")
		f.Meta()["goodplugin"] = name
	}
	return nil
}

func TestUsage(t *testing.T) {
	var stdout bytes.Buffer
	uconfig.UsageOutput = &stdout

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Usage should not panic, but did: %v", r)
		}
	}()

	// good plugin is used just so that we have more than
	// one tag/field that isn't pre-weighted in "usage".
	noopPlugin := &UselessPluginVisitor{}

	value := f.Config{}

	secretProvider := func(name string) (string, error) { return "top secret token", nil }

	c, err := uconfig.Classic(&value, nil, secret.New(secretProvider), noopPlugin)
	if err != nil {
		t.Fatal(err)
	}

	if size := stdout.Len(); size != 0 {
		t.Fatalf(
			"Expectedd nothing in UsageOutput before usage, got len: %v\n%s",
			size,
			stdout.String(),
		)
	}

	c.Usage()

	output := stdout.String()

	if diff := cmp.Diff(expectedUsageMessage, output); diff != "" {
		t.Error(diff)
	}
}
