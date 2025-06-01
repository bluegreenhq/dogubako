package log

import (
	"context"

	dogucontext "github.com/bluegreenhq/dogubako/context"
)

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, dogucontext.ContextKeyLogger, logger)
}

func ExtractLogger(ctx context.Context) Logger {
	return dogucontext.ExtractValue[Logger](ctx, dogucontext.ContextKeyLogger)
}
