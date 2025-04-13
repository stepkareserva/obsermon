package config

import "time"

type Config struct {
	Endpoint        string `env:"ADDRESS"`
	StoreIntervalS  int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func (c *Config) StoreInterval() time.Duration {
	return time.Duration(c.StoreIntervalS) * time.Second
}
