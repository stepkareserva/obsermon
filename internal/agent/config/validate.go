package config

import (
	"fmt"
	"net/url"
)

func Validate(c Config) error {
	_, err := url.ParseRequestURI(c.EndpointURL())
	if err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	if c.PollInterval() <= 0 {
		return fmt.Errorf("invalid poll interval %v", c.PollInterval())
	}
	if c.ReportInterval() <= 0 {
		return fmt.Errorf("invalid report interval %v", c.ReportInterval())
	}
	return nil
}
