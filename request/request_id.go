package request

import (
	"context"

	dogucontext "github.com/bluegreenhq/dogubako/context"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, dogucontext.ContextKeyRequestID, requestID)
}

func ExtractRequestID(ctx context.Context) string {
	return dogucontext.ExtractValue[string](ctx, dogucontext.ContextKeyRequestID)
}
