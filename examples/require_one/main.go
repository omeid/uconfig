package main

import (
	"encoding/json"
	"fmt"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/file"
)

// Config is our application config.
type Config struct {
	Address string `default:"localhost:8080" usage:"the address to bind to"`
	Mode    string `default:"start" flag:",command" usage:"start|stop"`
}

// files defines our mutually exclusive, required configuration files.
// uconfig will ensure exactly one of these files exists and parse it.
var files = file.RequireOne{
	{file.Relative("config.json"), json.Unmarshal},
	{file.Relative(".config.json"), json.Unmarshal},
	// Imagine supporting yaml too:
	// {"config.yaml", yaml.Unmarshal},
}

// conf sets up our manager with defaults, RequireOne file loader, env vars, and flags.
var conf = uconfig.Classic[Config](files)

func main() {
	// Run() will parse and check for errors, including printing
	// the usage message (with "Require One" files listed) if 0 or >1
	// of the required files exist, or if -h/--help is passed.
	c := conf.Run()

	fmt.Printf("Config loaded successfully! Address: %s, Mode: %s\n", c.Address, c.Mode)
}
