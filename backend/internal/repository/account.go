package repository

import (
	"context"
	"database/sql"
	"e_wallet/backend/domain"

	"github.com/doug-martin/goqu/v9"
)

type accountRepository struct {
	db *goqu.Database
}

func NewAccount(con *sql.DB) domain.AccountRepository {
	return &accountRepository{
		db: goqu.New("default", con),
	}
}

func (a accountRepository) FindByUserId(ctx context.Context, id int64) (account domain.Account, err error) {
	dataset := a.db.From("account").Where(goqu.Ex{
		"user_id": id,
	})

	_, err = dataset.ScanStructContext(ctx, &account)

	return
}

func (a accountRepository) FindByAccountNumber(ctx context.Context, accNumber string) (account domain.Account, err error) {
	dataset := a.db.From("account").Where(goqu.Ex{
		"account_number": accNumber,
	})

	_, err = dataset.ScanStructContext(ctx, &account)

	return
}

func (a accountRepository) Update(ctx context.Context, account *domain.Account) error {
	executor := a.db.Update("account").Where(goqu.Ex{
		"id": account.ID,
	}).Set(goqu.Record{
		"balance": account.Balance,
	}).Executor()

	_, err := executor.ExecContext(ctx)

	return err
}
