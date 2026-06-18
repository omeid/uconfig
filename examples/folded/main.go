package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/omeid/uconfig/plugins/fileset"

	"github.com/omeid/uconfig/examples/folded/database"
	"github.com/omeid/uconfig/examples/folded/redis"
)

// App is an application configuration within a fileset.
type App struct {
	Name string `usage:"the application name"`
	Port int    `usage:"port to bind the app to"`
	Meta struct {
		Version string   `usage:"app version tag"`
		Labels  []string `usage:"labels for the app"`
		Tags    struct {
			Name  string `usage:"tag name"`
			Value string `usage:"tag value"`
		}
	}
}

// Webhook is a webhook configuration.
type Webhook struct {
	URL    string   `usage:"destination url"`
	Events []string `usage:"list of events to subscribe to"`
}

// Config is our application config.
type Config struct {
	Hosts          []string                 `default:"localhost,localhost.local" usage:"the ip or domains to bind to"`
	RegionTimeouts map[string]time.Duration `default:"us:500ms,eu:1s,ap:1200ms" usage:"per-region request timeouts"`

	Redis    redis.Config
	Database database.Config

	Services struct {
		Apps []App `fileset:"apps"`
	}

	Webhooks []Webhook `fileset:"webhooks"`

	Mode string `default:"start" usage:"run|start|stop" flag:",command"`
}

var files = uconfig.Files{
	{file.Absolute("/etc/demo-app/config.json"), json.Unmarshal, true},
	{file.Relative("config.json"), json.Unmarshal, true},
	// or short form {"config.json", json.Unmarshal, true},
	// And, of course, you can add as many files
	// as you want, and they will be applied
	// in the given order.
}

var (
	fsPluginAbs      = fileset.New("apps", fileset.Absolute("/etc/app/*.yaml"), json.Unmarshal)
	fsPluginRel      = fileset.New("apps", fileset.Relative("apps.d/*.json"), json.Unmarshal)
	fsPluginWebhooks = fileset.New("webhooks", fileset.Relative("webhooks.d/*.json"), json.Unmarshal)
)

var conf = uconfig.Classic[Config](files, fsPluginAbs, fsPluginRel, fsPluginWebhooks)

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
