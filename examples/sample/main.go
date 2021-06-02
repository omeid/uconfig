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
	Hosts []string `default:"localhost,localhost.local"`

	Redis    redis.Config
	Database database.Config
}

func main() {

	conf := &Config{}

	files := uconfig.Files{
		"config.json": json.Unmarshal,
		// you can add more files if you like,
		// they will be applied in this order.
	}

	c, err := uconfig.Classic(&conf, files)
	if err != nil {
		c.Usage()
		os.Exit(1)
	}

	// use conf as you please.
	fmt.Printf("%#v", conf)

}
