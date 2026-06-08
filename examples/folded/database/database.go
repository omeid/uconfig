package database

// Config holds the database configurations.
type Config struct {
	Address  string `default:"localhost"`
	Port     string `default:"28015" env:".SERVICE_PORT"`
	Database string `default:"my-project" env:".DB"`
}
