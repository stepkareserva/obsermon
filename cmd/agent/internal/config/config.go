package config

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Endpoint       string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func ParseConfig() Config {
	endpoint := flag.String("a", "localhost:8080",
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22 (without protocol)")
	pollInterval := flag.Int("p", 2,
		"poll interval, in seconds, positive integer")
	reportInterval := flag.Int("r", 10,
		"report interval, in seconds, positive integer")

	flag.Parse()

	// !!! add prefix http to endpoint path, if not exists.
	// tests pass -a param without http prefix, but it's required
	if !strings.Contains(*endpoint, "://") {
		*endpoint = "http://" + *endpoint
	}

	return Config{
		Endpoint:       *endpoint,
		PollInterval:   time.Duration(*pollInterval) * time.Second,
		ReportInterval: time.Duration(*reportInterval) * time.Second,
	}
}

func (c *Config) Validate() error {
	u, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsuppoeted protocol scheme: %s", u.Scheme)
	}
	if c.PollInterval <= 0 {
		return fmt.Errorf("invalid poll interval %s", c.ReportInterval)
	}
	if c.ReportInterval <= 0 {
		return fmt.Errorf("invalid report interval %s", c.ReportInterval)
	}
	return nil
}
