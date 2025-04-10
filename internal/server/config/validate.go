package config

import (
	"fmt"
	"net"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}

	return nil
}
