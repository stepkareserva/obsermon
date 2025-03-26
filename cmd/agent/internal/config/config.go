package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	// endpoint address, without protocol.
	// to get endpoint URL call EndpointURL()
	endpoint string `env:"ADDRESS"`
	// pool interval in seconds.
	// to get duration call PollInterval()
	pollIntervalS int `env:"POLL_INTERVAL"`
	// pool interval in seconds.
	// to get duration call ReportInterval()
	reportIntervalS int `env:"REPORT_INTERVAL"`
}

func (c *Config) EndpointURL() string {
	return "http://" + c.endpoint
}

func (c *Config) PollInterval() time.Duration {
	return time.Duration(c.pollIntervalS) * time.Second
}

func (c *Config) ReportInterval() time.Duration {
	return time.Duration(c.reportIntervalS) * time.Second
}

func (c *Config) ParseCommandLine() error {

	defaultEndpoint := "localhost:8080"
	defaultPollInterval := 2
	defaultReportInterval := 10

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&c.endpoint, "a", defaultEndpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22 (without protocol)")
	fs.IntVar(&c.pollIntervalS, "p", defaultPollInterval,
		"poll interval, in seconds, positive integer")
	fs.IntVar(&c.reportIntervalS, "r", defaultReportInterval,
		"report interval, in seconds, positive integer")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

func (c *Config) ParseEnv() error {
	return env.Parse(c)
}

func (c *Config) Validate() error {
	_, err := url.ParseRequestURI(c.EndpointURL())
	if err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	if c.PollInterval() <= 0 {
		return fmt.Errorf("invalid poll interval %v", c.ReportInterval())
	}
	if c.ReportInterval() <= 0 {
		return fmt.Errorf("invalid report interval %v", c.ReportInterval())
	}
	return nil
}
