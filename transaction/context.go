package transaction

import (
	"context"

	dogucontext "github.com/bluegreenhq/dogubako/context"
)

func WithTransaction(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, dogucontext.ContextKeyTransaction, tx)
}

func ExtractTransaction(ctx context.Context) Transaction {
	return dogucontext.ExtractValue[Transaction](ctx, dogucontext.ContextKeyTransaction)
}
