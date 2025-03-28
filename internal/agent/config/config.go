package config

import (
	"time"
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
