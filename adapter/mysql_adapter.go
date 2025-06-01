package adapter

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/bluegreenhq/dogubako/config"
	dogusql "github.com/bluegreenhq/dogubako/sql"
	"github.com/bluegreenhq/dogubako/transaction"
)

type MySQLAdapter struct {
	db *sqlx.DB
}

func NewMySQLAdapter(c *config.MySQLConfig) (*MySQLAdapter, error) {
	db, err := sqlx.Open("mysql", c.DataSourceName())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)

	return &MySQLAdapter{
		db: db,
	}, nil
}

func (a MySQLAdapter) ExecTx(ctx context.Context, stmt dogusql.Statement) (sql.Result, error) {
	executor := a.getExecutor(ctx)

	result, err := executor.ExecContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (a MySQLAdapter) GetTx(ctx context.Context, dest any, stmt dogusql.Statement) error {
	executor := a.getExecutor(ctx)

	err := executor.GetContext(ctx, dest, stmt.Query(), stmt.Args()...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a MySQLAdapter) SelectTx(ctx context.Context, dest any, stmt dogusql.Statement) error {
	executor := a.getExecutor(ctx)

	err := executor.SelectContext(ctx, dest, stmt.Query(), stmt.Args()...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a MySQLAdapter) QueryxTx(ctx context.Context, stmt dogusql.Statement) (*sqlx.Rows, error) {
	executor := a.getExecutor(ctx)

	rows, err := executor.QueryxContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rows, nil
}

func (a MySQLAdapter) ExistsTx(ctx context.Context, stmt dogusql.Statement) (bool, error) {
	executor := a.getExecutor(ctx)

	rows, err := executor.QueryxContext(ctx, stmt.Query(), stmt.Args()...)
	if err != nil {
		return false, errors.WithStack(err)
	}

	defer func() {
		_ = rows.Close()
	}()

	return rows.Next(), nil
}

func (a MySQLAdapter) CountTx(ctx context.Context, stmt dogusql.Statement) (uint, error) {
	executor := a.getExecutor(ctx)

	var rec struct {
		Count uint `db:"count"`
	}

	err := executor.GetContext(ctx, &rec, stmt.Query(), stmt.Args()...)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return rec.Count, nil
}

func (a MySQLAdapter) Truncate(ctx context.Context, tableName string) error {
	query := fmt.Sprintf(`truncate %v`, tableName)

	_, err := a.ExecTx(ctx, dogusql.NewStatement(query, nil))
	if err != nil {
		return err
	}

	return nil
}

func (a MySQLAdapter) BeginTransaction() (transaction.Transaction, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return newTransaction(tx), nil
}

func (a MySQLAdapter) Close() error {
	err := a.db.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a MySQLAdapter) getExecutor(ctx context.Context) executor {
	tx := transaction.ExtractTransaction(ctx)
	if tx == nil {
		return a.db
	}

	sqlxTx, ok := tx.(*transactionImpl)
	if !ok {
		return a.db
	}

	return sqlxTx
}
