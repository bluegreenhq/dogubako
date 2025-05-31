package context

import (
	"context"

	"github.com/bluegreenhq/dogubako/model"
)

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
