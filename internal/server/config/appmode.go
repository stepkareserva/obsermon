package config

import "fmt"

type AppMode string

const (
	Quiet AppMode = "quiet"
	Dev   AppMode = "dev"
	Prod  AppMode = "prod"
)

func (a *AppMode) IsValid() bool {
	if a == nil {
		return false
	}
	switch *a {
	case Quiet, Dev, Prod:
		return true
	}
	return false
}

func (m *AppMode) String() string {
	if m == nil {
		return ""
	}
	return string(*m)
}

func (m *AppMode) Set(s string) error {
	switch s {
	case string(Quiet), string(Dev), string(Prod):
		*m = AppMode(s)
		return nil
	default:
		return fmt.Errorf("invalid app mode: %s", s)
	}
}
