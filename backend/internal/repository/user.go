package repository

import (
	"context"
	"database/sql"
	"e_wallet/backend/domain"

	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(con *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", con),
	}
}

func (u userRepository) FindByID(ctx context.Context, id int64) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.Ex{
		"id": id,
	})

	_, err = dataset.ScanStructContext(ctx, &user)
	return
}

func (u userRepository) FindByUserName(ctx context.Context, username string) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.Ex{
		"username": username,
	})

	_, err = dataset.ScanStructContext(ctx, &user)
	return
}
