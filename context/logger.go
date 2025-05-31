package context

import (
	"context"

	"github.com/bluegreenhq/dogubako/model"
)

func WithLogger(ctx context.Context, logger model.Logger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

func FromContext(ctx context.Context) model.Logger {
	if logger, ok := ctx.Value(ContextKeyLogger).(model.Logger); ok {
		return logger
	}

	return nil
}
