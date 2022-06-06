# uConfig [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/uconfig) [![Build Status](https://app.travis-ci.com/omeid/uconfig.svg?branch=master)](https://app.travis-ci.com/omeid/uconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/omeid/uconfig)](https://goreportcard.com/report/github.com/omeid/uconfig) [![Coverage](https://gocover.io/_badge/github.com/omeid/uconfig?update)](https://gocover.io/github.com/omeid/uconfig)


Lightweight, zero-dependency, and extendable configuration management.

uConfig is extremely light and extendable configuration management library with zero dependencies. Every aspect of configuration is provided through a _plugin_, which means you can have any combination of flags, environment variables, defaults, secret providers, Kubernetes Downward API, and what you want, and only what you want, through plugins.


uConfig takes the config schema as a struct decorated with tags, nesting is supported.

Supports all basic types, time.Duration, and any other type through `encoding.TextUnmarshaler` interface.
See the _[flat view](https://godoc.org/github.com/omeid/uconfig/flat)_ package for details.

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

  "github.com/omeid/uconfig"

  "$PROJECT/redis"
  "$PROJECT/database"
)

// Config is our application config.
type Config struct {

  // yes you can have slices.
  Hosts    []string `default:"localhost,localhost.local" usage:"the ip or domains to bind to"`

  Redis    redis.Config
  Database database.Config

}



func main() {

  conf := &Config{}

  confFiles := uconfig.Files{
     // filepath, Unmarshal func, Optional
    {"config.json", json.Unmarshal, true}
    // you can add more files if you like,
    // they will be applied in the given order.
  }

  c, err := uconfig.Classic(&conf, confFiles)
  if err != nil {
    // you could consider printing usage in case of error
    // like so:
    c.Usage()
    fmt.Println(err)
    os.Exit(1)
  }

  // use conf as you please.
  fmt.Printf("start with hosts set to: %#v\n", conf.Hosts)

}
```

Run this program with `-h` or `--help` and it would print out the usage:


```
flag provided but not defined: -x

Supported Fields:
FIELD                FLAG                  ENV                  DEFAULT                      USAGE
-----                -----                 -----                -------                      -----
Hosts                -hosts                HOSTS                localhost,localhost.local    the ip or domains to bind to
Redis.Address        -redis-address        REDIS_ADDRESS        redis-master
Redis.Port           -redis-port           REDIS_PORT           6379
Redis.Password       -redis-password       REDIS_PASSWORD
Redis.DB             -redis-db             REDIS_DB             0
Redis.Expire         -redis-expire         REDIS_EXPIRE         5s
Database.Address     -database-address     DATABASE_ADDRESS     localhost
Database.Port        -database-port        DATABASE_PORT        28015
Database.Database    -database-database    DATABASE_DATABASE    my-project
```


## Custom names:

Sometimes you might want to use a different env var, or flag name for backwards compatibility or other reasons, you have two options.

1. uconfig tag

You can change the name of a field as seen by unconfig.
This option supports the usual nesting prefixing.
See the port example below.

2. Plugin specific tags

Most plugins support controlling the field name as seen by that specific plugin.

This option does not support nesting prefixes.
See the Database field in the example below.


```go
package database

// Config holds the database configurations.
type Database struct {
  Address  string `default:"localhost"`
  Port     string `default:"28015" uconfig:"Service.Port"`
  Database string `default:"my-project" env:"DB_NAME" flag:"main-db-name"`
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

```
Supported Fields:
FIELD                    FLAG                      ENV                      DEFAULT                      USAGE
-----                    -----                     -----                    -------                      -----
Hosts                    -hosts                    HOSTS                    localhost,localhost.local    the ip or domains to bind to
Redis.Port               -redis-port               REDIS_PORT               6379
Redis.Password           -redis-password           REDIS_PASSWORD
Redis.DB                 -redis-db                 REDIS_DB                 0
Redis.Expire             -redis-expire             REDIS_EXPIRE             5s
Database.Address         -database-address         DATABASE_ADDRESS         localhost
Database.Service.Port    -database-service-port    DATABASE_SERVICE_PORT    28015
Database.Database        -main-db-db               DB_NAME                  my-project
exit status 1
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

import (

  "github.com/omeid/uconfig"
  "github.com/omeid/uconfig/plugins/secret"
)
// Creds is an example of a config struct that uses secret values.
type Creds struct {
  // by default, secret plugin will generate a name that is identical
  // to env plugin, SCREAM_SNAKE_CASE, so in this case it will be
  // APIKEY however, following the standard uConfig nesting rules
  // in Config struct below, it becomes CREDS_APIKEY.
  APIKey   string `secret:""`
  // or you can provide your own name, which will not be impacted
  // by nesting or the field name.
  APIToken string `secret:"API_TOKEN"`
}

type Config struct {
  Redis   Redis
  Creds   Creds
}


func main() {

  conf := &Config{}


  files := uconfig.Files{
    {"config.json", json.Unmarshal, false}
  }

   // secret.New accepts a function that maps a secret name to it's value.
   secretPlugin := secret.New(func(name string) (string, error) {
      // you're free to grab the secret based on the name from wherever
      // you please, aws secrets-manager, hashicorp vault, or wherever.
      value, ok := secretSource.Get(name)

      if !ok {
        return "", ErrSecretNotFound
      }

      return value, nil
  })

  // then you can use the secretPlugin with uConfig like any other plugin.
  // Lucky, uconfig.Classic allows passing more plugins, which means
  // you can simply do the following for flags, envs, files, and secrets!
  _, err := uconfig.Classic(&value, files, secretPlugin)
  if err != nil {
    t.Fatal(err)
  }

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

  conf := &YourConfigStruct{}

  // It will panic on error
  uconfig.Must(conf, defaults.New())

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
