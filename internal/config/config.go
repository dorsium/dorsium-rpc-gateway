package config

import (
	"os"
	"time"
)

// Config holds application configuration.
type Config struct {
	Address        string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	NodeRPC        string
	Version        string
	Mode           string
	DisableMetrics bool
	AdminToken     string
}

// New creates a Config with default values.
func New() *Config {
	c := &Config{
		Address:      ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		NodeRPC:      "http://localhost:26657",
		Version:      "dev",
		Mode:         "production",
		AdminToken:   "changeme",
	}
	if addr := os.Getenv("ADDRESS"); addr != "" {
		c.Address = addr
	}
	if rpc := os.Getenv("NODE_RPC"); rpc != "" {
		c.NodeRPC = rpc
	}
	if v := os.Getenv("APP_VERSION"); v != "" {
		c.Version = v
	}
	if m := os.Getenv("APP_MODE"); m != "" {
		c.Mode = m
	}
	if t := os.Getenv("ADMIN_TOKEN"); t != "" {
		c.AdminToken = t
	}
	c.DisableMetrics = os.Getenv("DISABLE_METRICS") == "true"
	return c
}
