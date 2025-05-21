package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

func LoadConfig() (*Config, error) {
	cfg := defaultConfig()
	if err := readCLIParams(cfg); err != nil {
		return nil, err
	}
	if err := readEnvParams(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Endpoint:        "localhost:8080",
		PollIntervalS:   2,
		ReportIntervalS: 10,
	}
}

func readCLIParams(c *Config) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80,\n"+
			"localhost:22 (without protocol)")
	fs.IntVar(&c.PollIntervalS, "p", c.PollIntervalS,
		"poll (local metrics update) interval, in seconds,\n"+
			"positive integer")
	fs.IntVar(&c.ReportIntervalS, "r", c.ReportIntervalS,
		"report (send metrics to server) interval, in seconds,\n"+
			"positive integer")
	fs.StringVar(&c.ReportSignKey, "k", c.ReportSignKey,
		"report sign key, string")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

func readEnvParams(c *Config) error {
	return env.Parse(c)
}
