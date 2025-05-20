package context

import (
	"context"
	"time"

	"github.com/bluegreenhq/dogubako/model"
)

type ContextKey string

const (
	// userinterface layer
	ContextKeyRequestID   ContextKey = "requestId"
	ContextKeyRequestTime ContextKey = "requestTime"
	// application layer
	ContextKeyTransaction ContextKey = "transaction"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

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

func ExtractTransaction(ctx context.Context) model.Transaction {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(ContextKeyTransaction)
	if val == nil {
		return nil
	}

	tx, ok := val.(model.Transaction)
	if !ok {
		return nil
	}

	return tx
}

func WithTransaction(ctx context.Context, tx model.Transaction) context.Context {
	return context.WithValue(ctx, ContextKeyTransaction, tx)
}
