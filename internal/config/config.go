package config

import "time"

// Config holds application configuration.
type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	NodeRPC      string
}

// New creates a Config with default values.
func New() *Config {
	return &Config{
		Address:      ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		NodeRPC:      "http://localhost:26657",
	}
}
