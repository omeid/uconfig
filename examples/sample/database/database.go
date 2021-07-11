package database

// Config holds the database configurations.
type Config struct {
	Address  string `default:"localhost"`
	Port     string `default:"28015"`
	Database string `default:"my-project"`
}
