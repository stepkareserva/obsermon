package config

import (
	"fmt"
	"net"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	if c.StoreInterval() < 0 {
		return fmt.Errorf("invalid poll interval %v", c.StoreInterval())
	}
	if !c.Mode.IsValid() {
		return fmt.Errorf("invalid app mode %v", c.Mode)
	}
	// ? maybe some methods exists for this check?
	// if err := checkAccessRights(c.FileStoragePath); err != nil {
	//	return fmt.Errorf("invalid storage file: %w", err)
	//}

	return nil
}
