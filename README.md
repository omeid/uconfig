# uConfig [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/uconfig)  [![Build Status](https://travis-ci.org/omeid/uconfig.svg?branch=master)](https://travis-ci.org/omeid/uconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/omeid/uconfig)](https://goreportcard.com/report/github.com/omeid/uconfig)

uConfig is an unopinionated, extendable and plugable configuration management.

Every aspect of configuration is provided through a plugin, which means you can have any combination of flags, environment variables, defaults, Kubernetes Downward API, and what you want, through plugins.


uConfig takes the config schema as a struct decorated with tags, nesting is supported.

Supports all basic types, time.Duration, time.Time, and you any other type through `encoding.TextUnmarshaler` interface.
See the _[flat view](https://godoc.org/github.com/omeid/uconfig/flat)_ package for details.

## Example Configuration: 

```go
package database
// Config holds the database configurations.
type Config struct {
  Address  string `default:"localhost" env:"DATABASE_HOST"`
  Port     string `default:"28015" env:"DATABASE_SERVICE_PORT"`
  Database string `default:"my-project"`
}
```
```go
package redis
// Config describes the requirement for redis client.
type Config struct {
  Address  string `default:"redis-master" env:"REDIS_HOST"`
  Port     string `default:"6379" env:"REDIS_SERVICE_PORT"`
  Password string `default:""`
  DB       int    `default:"0"`
}
```

```go
package main
// Config is our distribution configs as required for services and clients.
type Config struct {

  Redis   redis.Config
  Database database.Config
}

```


## Plugins

uConfig supports two kind of plugins, Walkers and Visitors.

### Walkers 

Walkers are used for configuration plugins that take the whole config struct and unmarshal the underlying content into the config struct.
Plugins that load the configuration from files are good candidates for this.

```go
// Walker is the interface for plugins that take the whole config, like file loaders.
type Walker interface {
  Plugin

  Walk(interface{}) error
}
```


### Visitors

Visitors get a _[flat view](https://godoc.org/github.com/omeid/uconfig/flat)_ of the configuration struct, which is a flat view of the structs regardless of nesting level, for more details see the [flat](https://godoc.org/github.com/omeid/uconfig/flat) package documentation.

Plugins that load the configurations from flat structures (e.g flags, environment variables, default tags) are good candidts for this type of plugin.


```go
// Visitor is the interface for plugins that require a flat view of the config, like flags, env vars
type Visitor interface {
  Plugin

  Visit(flat.Fields) error
}

```


## File Plugin

File plugin is a walker plugin that loads configuration files of different formats by way of accepting an Unmarshaler function that follows the standard unmarshal function of type `func(src []byte, v interface{}) error`; this allows you to use `encoding/json` and other encoders that follow the same interface. 

Following is some common unmarshalers that follow the standard unmarshaler function:

* JSON: `encoding/json`
* TOML: `github.com/BurntSushi/toml`
* YAML: `gopkg.in/yaml.v2`
  * Note: YAML unmarshaller doesn't appear to handle embedded structs as cleanly as some unmarshallers; you may need to nest the embedded struct's options in your YAML file (see `version` in the example below)


## Example

The following example uses `uconfig.Classic` to create a uConfig manager which processes defaults, optionally any config files, environment variables, and flags; in that order.
In this example, we're using a single YAML config file, but you can specify multiple files (each with its own unmarshaller) in the `uconfig.Files` map if required.

```yaml
# path/to/config.yaml
# YAML unmarshaller doesn't appear to handle flattened embedded structs,
# so 'version' needs to sit under 'anon' to map correctly
anon:
  version: '0.2'
gohard: true
redis:
  host: redis-host
  port: 6379
rethink:
  db: base
  host:
    address: rethink-cluster
    port: '28015'
```

```go
// main.go
package main

import (
  "os"

  "gopkg.in/yaml.v2"

  "github.com/omeid/uconfig"
)

type Anon struct {
  Version string `default:"0.0.1" env:"APP_VERSION"`
}

type Host struct {
  Address string `default:"localhost" env:"RETHINKDB_HOST"`
  Port    string `default:"28015" env:"RETHINKDB_PORT"`
}

type RethinkConfig struct {
  Host Host
  Db   string `default:"my-project"`
}

type Redis struct {
  Host string `default:"redis-master" env:"REDIS_HOST"`
  Port int    `default:"6379" env:"REDIS_SERVICE_PORT"`
}

type YourConfig struct {
  Anon
  GoHard  bool
  Redis   Redis
  Rethink RethinkConfig
}

func main() {

  conf := &YourConfig{}

  // Simply
  c, err := uconfig.Classic(conf, uconfig.Files{
    "path/to/config.yaml": yaml.Unmarshal,
  })
  if err != nil {
    c.Usage()
    os.Exit(1)
  }
  // Use your config here as you please.
}

```

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
