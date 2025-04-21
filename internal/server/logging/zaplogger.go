package logging

import (
	"fmt"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(m config.AppMode) (*zap.Logger, error) {
	if m == config.Quiet {
		return zap.NewNop(), nil
	}

	var cfg zap.Config
	var options []zap.Option

	switch m {
	case config.Dev:
		cfg = zap.NewDevelopmentConfig()
		options = append(options, zap.WithCaller(false))
	case config.Prod:
		cfg = zap.NewProductionConfig()
		options = append(options, zap.WithCaller(false))
	default:
		return nil, fmt.Errorf("invalid zap log level")
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build(options...)
}
