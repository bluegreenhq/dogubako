package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type executor interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}
