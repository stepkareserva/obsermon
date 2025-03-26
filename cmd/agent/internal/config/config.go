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
	Endpoint string `env:"ADDRESS"`
	// pool interval in seconds.
	// to get duration call PollInterval()
	PollIntervalS int `env:"POLL_INTERVAL"`
	// pool interval in seconds.
	// to get duration call ReportInterval()
	ReportIntervalS int `env:"REPORT_INTERVAL"`
}

func (c *Config) EndpointURL() string {
	return "http://" + c.Endpoint
}

func (c *Config) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalS) * time.Second
}

func (c *Config) ReportInterval() time.Duration {
	return time.Duration(c.ReportIntervalS) * time.Second
}

func (c *Config) ParseCommandLine() error {

	defaultEndpoint := "localhost:8080"
	defaultPollInterval := 2
	defaultReportInterval := 10

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&c.Endpoint, "a", defaultEndpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80,\n"+
			"localhost:22 (without protocol)")
	fs.IntVar(&c.PollIntervalS, "p", defaultPollInterval,
		"poll (local metrics update) interval, in seconds,\n"+
			"positive integer")
	fs.IntVar(&c.ReportIntervalS, "r", defaultReportInterval,
		"report (send metrics to server) interval, in seconds,\n"+
			"positive integer")

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
