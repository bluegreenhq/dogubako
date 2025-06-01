package request

import (
	"context"
	"time"

	dogucontext "github.com/bluegreenhq/dogubako/context"
)

func WithRequestTime(ctx context.Context, requestTime time.Time) context.Context {
	return context.WithValue(ctx, dogucontext.ContextKeyRequestTime, requestTime)
}

func ExtractRequestTime(ctx context.Context) time.Time {
	return dogucontext.ExtractValue[time.Time](ctx, dogucontext.ContextKeyRequestTime)
}
