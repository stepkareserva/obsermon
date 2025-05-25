package config

import (
	"flag"
	"os"
	"path"
	"runtime"

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
		StoreIntervalS:  300,
		FileStoragePath: defaultStoragePath(),
		Restore:         false,
		DBConnection:    "",
		ReportSignKey:   "",
		Mode:            Prod,
	}
}

func defaultStoragePath() string {
	// dir for apps data - depends from os
	var appDataDir string
	var err error
	switch runtime.GOOS {
	case "windows", "darwin":
		appDataDir, err = os.UserConfigDir()
	default:
		appDataDir, err = os.UserHomeDir()
	}

	if err != nil {
		appDataDir = ""
	}

	appDataDir = path.Join(appDataDir, "obsermon")
	if err := os.MkdirAll(appDataDir, os.ModePerm); err != nil {
		appDataDir = ""
	}

	return path.Join(appDataDir, "storage.json")
}

func readCLIParams(c *Config) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")

	fs.IntVar(&c.StoreIntervalS, "i", c.StoreIntervalS,
		"server state storing interval, s, 0 for sync storing")

	fs.StringVar(&c.FileStoragePath, "f", c.FileStoragePath,
		"path to server state storage file")

	fs.BoolVar(&c.Restore, "r", c.Restore,
		"restore server state from storage file")

	fs.StringVar(&c.DBConnection, "d", c.DBConnection,
		"database connection string")

	fs.StringVar(&c.ReportSignKey, "k", c.ReportSignKey,
		"reports and requests sing key, string")

	fs.Var(&c.Mode, "m",
		"app mode, quiet/dev/prod")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

func readEnvParams(c *Config) error {
	return env.Parse(c)
}
