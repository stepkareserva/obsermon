package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(level Level) (*zap.Logger, error) {
	if level == LevelNoop {
		return zap.NewNop(), nil
	}

	var cfg zap.Config
	var options []zap.Option

	switch level {
	case LevelDev:
		cfg = zap.NewDevelopmentConfig()
		options = append(options, zap.WithCaller(false))
	case LevelProd:
		cfg = zap.NewProductionConfig()
		options = append(options, zap.WithCaller(false))
	default:
		return nil, fmt.Errorf("invalid zap log level")
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build(options...)
}
