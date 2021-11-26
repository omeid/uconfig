package uconfig_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/internal/f"
)

const expectedUsageMessage = `
Supported Fields:
FIELD                   FLAG                     ENV                     DEFAULT    USAGE
-----                   -----                    -----                   -------    -----
Version                 -version                 VERSION                            
GoHard                  -gohard                  GOHARD                             
Redis.Host              -redis-host              REDIS_HOST                         
Redis.Port              -redis-port              REDIS_PORT                         
Rethink.Host.Address    -rethink-host-address    RETHINK_HOST_ADDRESS               
Rethink.Host.Port       -rethink-host-port       RETHINK_HOST_PORT                  
Rethink.Db              -rethink-db              RETHINK_DB              primary    main database used by our application
`

func TestUsage(t *testing.T) {
	var stdout bytes.Buffer
	uconfig.UsageOutput = &stdout

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Usage should not panic, but did: %v", r)
		}
	}()

	conf := f.Config{}
	c, err := uconfig.Classic(&conf, nil)
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
