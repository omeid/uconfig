package redis

import "time"

// Config describes the requirement for redis client.
type Config struct {
	Address  string        `default:"redis-master"`
	Port     string        `default:"6379"`
	Password string        `secret:""`
	DB       int           `default:"0"`
	Expire   time.Duration `default:"5s"`
}
