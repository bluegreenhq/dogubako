package context

import (
	"context"
	"time"
)

func ExtractRequestTime(ctx context.Context) time.Time {
	if ctx == nil {
		return time.Time{}
	}

	val := ctx.Value(ContextKeyRequestTime)
	if val == nil {
		return time.Time{}
	}

	requestTime, ok := val.(time.Time)
	if !ok {
		return time.Time{}
	}

	return requestTime
}

func WithRequestTime(ctx context.Context, requestTime time.Time) context.Context {
	return context.WithValue(ctx, ContextKeyRequestTime, requestTime)
}
