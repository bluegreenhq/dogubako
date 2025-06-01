package adapter

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/bluegreenhq/dogubako/transaction"
)

type transactionImpl struct {
	*sqlx.Tx
	id string
}

var _ transaction.Transaction = (*transactionImpl)(nil)

func (t *transactionImpl) ID() string {
	return t.id
}

func newTransaction(tx *sqlx.Tx) *transactionImpl {
	return &transactionImpl{
		Tx: tx,
		id: uuid.NewString(),
	}
}
