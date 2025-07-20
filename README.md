# uConfig [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/uconfig) [![Build Status](https://app.travis-ci.com/omeid/uconfig.svg?branch=master)](https://app.travis-ci.com/omeid/uconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/omeid/uconfig)](https://goreportcard.com/report/github.com/omeid/uconfig) [![Coverage](https://gocover.io/_badge/github.com/omeid/uconfig?update)](https://gocover.io/github.com/omeid/uconfig)


Lightweight, zero-dependency, and extendable configuration management.

uConfig is extremely light and extendable configuration management library with zero dependencies. Every aspect of configuration is provided through a _plugin_, which means you can have any combination of flags, environment variables, defaults, secret providers, Kubernetes Downward API, and any combination of configuration files and formats including json, toml, cue, or just about anything you want, and only what you want, through plugins.


To use uConfig, you simply define the configuration struct for your services and application, and uConfig does all the heavy-lifting. It just works.

## Example Configuration:

```go
package database
// Config holds the database configurations.
type Config struct {
  Address  string `default:"localhost"`
  Port     string `default:"28015"`
  Database string `default:"my-project"`
}
```

```go
package redis
// Config describes the requirement for redis client.
type Config struct {
  Address  string        `default:"redis-master"`
  Port     string        `default:"6379"`
  Password string        `secret:""`
  DB       int           `default:"0"`
  Expire   time.Duration `default:"5s"`
}
```

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/omeid/uconfig"

	"github.com/omeid/uconfig/examples/sample/database"
	"github.com/omeid/uconfig/examples/sample/redis"
)


// Config is our application config.
type Config struct {
	// yes you can have slices.
	Hosts []string `default:"localhost,localhost.local" usage:"the ip or domains to bind to"`

	Redis    redis.Config
	Database database.Config

	// the flags plugin allows capturing a single Command after the flags.
	// so you can run myprogram -flag=value -s -blah=bleh stop|start|stop and so on.
	Mode string `default:"start" flag:",command" usage:"run|start|stop"`
}

var files = uconfig.Files{
	{Path: "/etc/demo-app/config.json", Unmarshal: json.Unmarshal, Optional: true},
	{Path: "config.json", Unmarshal: json.Unmarshal, Optional: true},
	// or short form {"config.json", json.Unmarshal, true},
	// And, of course, you can of course add as many files
	// as you want, and they will be applied
	// in the given order.
}

var conf = uconfig.Classic[Config](files)

func main() {
	conf := conf.Run()
	// use conf as you please.
	// let's pretty print it as JSON for example:
	configAsJson, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(string(configAsJson))
}
```

Now lets run our program:

```sh
$ go run main.go -h
Usage:
    main [flags] [command]

Configurations:
FIELD                FLAG                  ENV                 DEFAULT                      USAGE
-----                -----                 -----               -------                      -----
Hosts                -hosts                HOSTS               localhost,localhost.local    the ip or domains to bind to
Redis.Address        -redis-address        REDIS_ADDRESS       redis-master                 
Redis.Port           -redis-port           REDIS_PORT          6379                         
Redis.Password       -redis-password       REDIS_PASSWORD                                   
Redis.DB             -redis-db             REDIS_DB            0                            
Redis.Expire         -redis-expire         REDIS_EXPIRE        5s                           
Database.Address     -database-address     DATABASE_ADDRESS    localhost                    
Database.Port        -database-port        SERVICE_PORT        28015                        
Database.Database    -database-database    DB                  my-project                   
Mode                 [command]             MODE                start                        run|start|stop

Configuration Files:
    /etc/demo-app/config.json
    config.json

```
```sh
$ go run main.go 
```
```json
{
 "Hosts": [
  "localhost",
  "localhost.local"
 ],
 "Redis": {
  "Address": "redis-master",
  "Port": "6379",
  "Password": "",
  "DB": 0,
  "Expire": 5000000000
 },
 "Database": {
  "Address": "localhost",
  "Port": "28015",
  "Database": "my-project"
 },
 "Mode": "start"
}

```

uConfig supports all basic types, time.Duration, slices, and any other type through `encoding.TextUnmarshaler` interface.
See the _[flat view](https://godoc.org/github.com/omeid/uconfig/flat)_ package for details.

## Custom names:

Sometimes you might want to use a different env var, or flag name for backwards compatibility or other reasons, you have two options.

1. uconfig tag

You can change the name of a field as seen by `uconfig`.

Please note that this flag only works for walker plugins (flags, env, anything flat) and for Visitor plugins (file, stream, et al) you will need to use encoder specific tags like `json:"field_name"` and so on.


2. Plugin specific tags

Most plugins support controlling the field name as seen by that specific plugin. For example `env:"DB_NAME"`.


For both type of tags, you can prefix them with `.` to rename the field only at the struct level.
See the `Service.Port` and `DB_NAME` examples below.

```go
package database

// Config holds the database configurations.
type Database struct {
  Address  string `default:"localhost"`
  Port     string `default:"28015" uconfig:".Service.Port"` // field level rename.
  Database string `default:"my-project" env:"DB_NAME" flag:"main-db-name"` // plugin specific rename.
}
```


```go
package main

// Config is our application config.
type Config struct {

  // yes you can have slices.
  Hosts    []string `default:"localhost,localhost.local"`

  Redis    redis.Config
  Database database.Config
}
```


Which should give you the following settings:

```sh
$ go run main.go -h
Usage:
    main [flags] [command]

Configurations:
FIELD                  FLAG                    ENV                    DEFAULT                      USAGE
-----                  -----                   -----                  -------                      -----
Hosts                  -hosts                  HOSTS                  localhost,localhost.local    the ip or domains to bind to
Redis.Address          -redis-address          REDIS_ADDRESS          redis-master                 
Redis.Port             -redis-port             REDIS_PORT             6379                         
Redis.Password         -redis-password         REDIS_PASSWORD                                      
Redis.DB               -redis-db               REDIS_DB               0                            
Redis.Expire           -redis-expire           REDIS_EXPIRE           5s                           
Database.Address       -database-address       DATABASE_ADDRESS       localhost                    
Database.Database      -main-db-db             DB_NAME                my-project
Database.Service.Port  -database-service-port  DATABASE_SERVICE_PORT  28015
```


For file based plugins, you will need to use the appropriate tags as used by your encoder of choice. For example:

```go
package users

// Config holds the database configurations.
type Config struct {
  Host string `json:"bind_addr"`
}
```

## Secrets Plugin
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/uconfig/plugins/secret)

The secret provider allows you to grab the value of a config from anywhere you want. You simply need to implement the `func(name string) (value string)` function and pass it to the secrets plugin.

Unlike most other plugins, secret requires explicit `secret:""` tag, this is because only specific config values like passwords and api keys come from a secret provider, compared to the rest of the config which can be set in various ways.

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/omeid/uconfig"
    "github.com/omeid/uconfig/plugins/secret"

    "github.com/omeid/uconfig/examples/secrets/secretsource"
)

// Creds is an example of a config struct that uses secret values.
type Creds struct {
    // by default, secret plugin will generate a name that is identical
    // to env plugin, SCREAM_SNAKE_CASE, so in this case it will be
    // APIKEY however, following the standard uConfig nesting rules
    // in Config struct below, it becomes CREDS_APIKEY.
    APIKey string `secret:""`
    // or you can provide your own name, which will not be impacted
    // by nesting or the field name.
    APIToken string `secret:"API_TOKEN"`
}

type Config struct {
    Creds Creds
}

var secrets = secret.New(func(name string) (string, error) {
    // you're free to grab the secret based on the name from wherever
    // you please, aws secrets-manager, hashicorp vault, or wherever.
    value, ok := secretsource.Get(name)
    if !ok {
        return "", secret.ErrSecretNotFound
    }

    return value, nil
})

func main() {
    // then you can use the secretPlugin with uConfig like any other plugin.
    // Lucky, uconfig.Classic allows passing more plugins, which means
    // you can simply do the following for flags, envs, files, and secrets!
    conf := uconfig.Classic[Config](nil, secrets).Run()

    fmt.Printf("we got an API Key: %s\n", conf.Creds.APIKey)
}
```


## Tests

For tests, you may consider the `Must` function to set the defaults, like so
```go
package something

import (
  "testing"

  "github.com/omeid/uconfig"
  "github.com/omeid/uconfig/defaults"
)

func TestSomething(t *testing.T) error {

  // It will panic on error
    conf := uconfig.Must[Conf](defaults.New())

  // Use your conf as you please.
}

```

See the Classic source for how to compose plugins.
For more details, see the [godoc](https://godoc.org/github.com/omeid/uconfig).

## Extending uConfig:

uConfig provides a plugin mechanism for adding new sources of configuration.
There are two kind of plugins, Walkers and Visitors.

To implement your own, see the examples.


### Visitors

Visitors get a _[flat view](https://godoc.org/github.com/omeid/uconfig/flat)_ of the configuration struct, which is a flat view of the structs regardless of nesting level, for more details see the [flat](https://godoc.org/github.com/omeid/uconfig/flat) package documentation.

Plugins that load the configurations from flat structures (e.g flags, environment variables, default tags) are good candidates for this type of plugin.
See [env plugin](plugins/env/env.go) for an example.

### Walkers

Walkers are used for configuration plugins that take the whole config struct and unmarshal the underlying content into the config struct.
Plugins that load the configuration from files are good candidates for this.

See [file plugin](plugins/file/file.go) for an example.
