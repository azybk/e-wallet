package repository

import (
	"context"
	"database/sql"
	"e_wallet/backend/domain"

	"github.com/doug-martin/goqu/v9"
)

type TransactionRepository struct {
	db *goqu.Database
}

func NewTransaction(con *sql.DB) domain.TransactionRepository {
	return &TransactionRepository{
		db: goqu.New("default", con),
	}
}

func (t TransactionRepository) Insert(ctx context.Context, transaction *domain.Transaction) error {
	executor := t.db.Insert("transaction").Rows(goqu.Record{
		"account_id":           transaction.AccountId,
		"sof_number":           transaction.SofNumber,
		"dof_number":           transaction.DofNumber,
		"transaction_type":     transaction.TransactionType,
		"amount":               transaction.Amount,
		"transaction_datetime": transaction.TransactionDatetime,
	}).Returning("id").Executor()

	_, err := executor.ScanStructContext(ctx, transaction)

	return err
}
