package context

import "context"

type ContextKey string

const (
	// userinterface layer
	ContextKeyRequestID   ContextKey = "requestId"
	ContextKeyRequestTime ContextKey = "requestTime"
	// application layer
	ContextKeyTransaction ContextKey = "transaction"
	ContextKeyLogger      ContextKey = "logger"
)

func ExtractValue[T any](ctx context.Context, key ContextKey) T {
	if value, ok := ctx.Value(key).(T); ok {
		return value
	}

	var zero T

	return zero
}
