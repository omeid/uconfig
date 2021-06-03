# uConfig with YAML


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
    { "path/to/config.yaml", yaml.Unmarshal },
  })

  // or alternatively, using your own combination of plugins
  // see uconfig.Classic function for an example.

  if err != nil {
    c.Usage()
    os.Exit(1)
  }
  // Use your config here as you please.
}

```
