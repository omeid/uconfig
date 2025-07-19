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
