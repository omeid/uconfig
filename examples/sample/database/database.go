package database

// Config holds the database configurations.
type Config struct {
	Address  string `default:"localhost" env:"DATABASE_HOST"`
	Port     string `default:"28015" env:"DATABASE_SERVICE_PORT"`
	Database string `default:"my-project"`
}
