# uConfig [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/uconfig)  [![Build Status](https://travis-ci.org/omeid/uconfig.svg?branch=master)](https://travis-ci.org/omeid/uconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/omeid/uconfig)](https://goreportcard.com/report/github.com/omeid/uconfig)

uConfig is an unopinionated, extendable and plugable configuration management.

Every aspect of configuration is provided through a plugin, which means you can have any combination of flags, environment variables, defaults, Kubernetes Downward API, and you want it, through plugins.


uConfig takes the config schema as a struct decorated with tags, nesting is supported.


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


## Example



```go
package main

import (
	"log"
	"os"

	"github.com/omeid/uconfig"
)

func main() error {

	conf := &YourConfigStruct{}

	// Simply
	c, err := uconfig.Classic(conf, uconfig.Files{
		"path/to/config.json": json.Unmarshal,
		"path/to/config.toml": toml.Unmarshal,
	})
	if err != nil {
		c.Usage()
		os.Exit(1)
	}
	// User your config here as you please.
}

```

See the Classic source for how to compose plugins.  
For more details, see the [godoc](https://godoc.org/github.com/omeid/uconfig).
