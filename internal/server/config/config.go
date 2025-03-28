package config

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Endpoint string `env:"ADDRESS"`
}

func (c *Config) ParseCommandLine() error {
	defaultEndpoint := "localhost:8080"

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&c.Endpoint, "a", defaultEndpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

func (c *Config) ParseEnv() error {
	return env.Parse(c)
}

func (c *Config) Validate() error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}

	return nil
}
