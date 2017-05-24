# uConfig [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/uconfig)  [![Build Status](https://travis-ci.org/omeid/uconfig.svg?branch=master)](https://travis-ci.org/omeid/uconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/omeid/uconfig)](https://goreportcard.com/report/github.com/omeid/uconfig)

uConfig is an opinionated* extendable and plugable configuration management.

Every aspect of configuration is provided through a plugin, which means you can have any combination of flags, environment variables, defaults, Kubernetes Downward API, and you want it, through plugins.


uConfig takes the config schema as a struct decorated with tags, nesting is supported.


## Example Configuration: 

```go
// package database 
// Config holds the database configurations.
type Config struct {
  Address  string `default:"localhost" env:"DATABASE_HOST"`
  Port     string `default:"28015" env:"DATABASE_SERVICE_PORT"`
  Database string `default:"my-project"`
}

// package redis
// Config describes the requirement for redis client.
type Config struct {
  Address  string `default:"redis-master" env:"REDIS_HOST"`
  Port     string `default:"6379" env:"REDIS_SERVICE_PORT"`
  Password string `default:""`
  DB       int    `default:"0"`
}

// package main
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

```
// Walker is the interface for plugins that take the whole config, like file loaders.
type Walker interface {
  Plugin

  Walk(interface{}) error
}
```


### Visitors

Visitors get a _flatview_ of the configuration struct, which is a flat view of the structs regardless of nesting level, for more details see the flatview package documentation.

Plugins that load the configurations from flat structures (e.g flags, environment variables, default tags) are good candidts for this type of plugin.


```go
// WalkerFunc is a helper type that turns a Walk function into a Walker.
// type WalkerFunc func(interface{}) error

// Walk implements Walker for WalkerFunc
// func (wf WalkerFunc) Walk(conf interface{}) error { return wf(conf) }

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
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/flag"
)

func main() error {

	conf := &YourConfigStruct{}

	// Simply
	c, err := uconfig.CLassic(conf)
	if err != nil {
		c.Usage()
		os.Exit(1)
	}

	// Or if you want more control over what plugins are loaded,
	// the Classic helper is equivalent to this:

	// Start Classic
	c, err := uconfig.New(conf,
		// We use plugins to add support for loading config
		// from different source. The order is important!

		// Loads the default values from the default tag.
		defaults.New(),
		// Loads the configurations from the env vars.
		env.New(),
		// Loads the configurations from the flags.
		flag.Standard(),
	)

	if err != nil {
		log.Fatal(err)
	}
	err := c.Parse()
	if err != nil {
		c.Usage()
		log.Fatal(err)
	}
	// End Classic

	// User your config here as you please.
}

```

For more details, see the [godoc](https://godoc.org/github.com/omeid/uconfig).
