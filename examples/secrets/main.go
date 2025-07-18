package main

import (
	"encoding/json"
	"fmt"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/examples/secrets/secretsource"
	"github.com/omeid/uconfig/plugins/secret"
)

// Creds is an example of a config struct that uses secret values.
type Creds struct {
	// by default, secret plugin will generate a name that is identical
	// to env plugin, SCREAM_SNAKE_CASE, so in this case it will be
	// APIKEY however, following the standard uConfig nesting rules
	// in Config struct below, it becomes CREDS_APIKEY.
	APIKey string `secret:""`
	// or you can provide your own name, which will not be impacted
	// by nesting or the field name.
	APIToken string `secret:"API_TOKEN"`
}

type Config struct {
	Creds Creds
}

var files = uconfig.Files{
	{Path: "config.json", Unmarshal: json.Unmarshal, Optional: true},
	// or short form: {"config.json", json.Unmarshal, true},
}

var secrets = secret.New(func(name string) (string, error) {
	// you're free to grab the secret based on the name from wherever
	// you please, aws secrets-manager, hashicorp vault, or wherever.
	value, ok := secretsource.Get(name)
	if !ok {
		return "", secret.ErrSecretNotFound
	}

	return value, nil
})

func main() {
	// then you can use the secretPlugin with uConfig like any other plugin.
	// Lucky, uconfig.Classic allows passing more plugins, which means
	// you can simply do the following for flags, envs, files, and secrets!
	conf := uconfig.Classic[Config](files, secrets).Run()

	fmt.Printf("we got an API Key: %s\n", conf.Creds.APIKey)
}
