package config

import (
	"flag"
	"fmt"
	"net"
)

type Config struct {
	Endpoint string
}

func ParseConfig() Config {
	var config Config

	flag.StringVar(&config.Endpoint, "a", "localhost:8080",
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")

	flag.Parse()
	return config
}

func (c *Config) Validate() error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}

	return nil
}
