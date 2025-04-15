package logging

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(contextKey{}).(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}
