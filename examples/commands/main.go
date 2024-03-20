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
}

func main() {

	files := uconfig.Files{
		{"config.json", json.Unmarshal, true},
	}

	mainFunc := func(conf Config) error {
		// use conf as you please.
		// let's pretty print it as JSON for example:
		configAsJson, err := json.MarshalIndent(conf, "", " ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print(string(configAsJson))
		return nil
	}

	destroyFunc := func(conf struct{}) error {
		fmt.Println("Running destroy with no config")
		return nil
	}

	err := uconfig.Commands(
		uconfig.ClassicCommand("", mainFunc, files),
		uconfig.ClassicCommand("destroy", destroyFunc, nil),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
