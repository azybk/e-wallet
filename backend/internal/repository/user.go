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

func (u userRepository) Insert(ctx context.Context, user *domain.User) error {
	executors := u.db.Insert("users").Rows(goqu.Record{
		"full_name": user.FullName,
		"phone":     user.Phone,
		"username":  user.Username,
		"password":  user.Password,
		"email":     user.Email,
	}).Returning("id").Executor()

	_, err := executors.ScanStructContext(ctx, user)

	return err
}

func (u userRepository) Update(ctx context.Context, user *domain.User) error {

	user.EmailVerifiedAtDB = sql.NullTime{
		Time:  user.EmailVerifiedAt,
		Valid: true,
	}

	executor := u.db.Update("users").Set(goqu.Record{
		"full_name":         user.FullName,
		"phone":             user.Phone,
		"username":          user.Username,
		"password":          user.Password,
		"email":             user.Email,
		"email_verified_at": user.EmailVerifiedAt,
	}).Executor()

	_, err := executor.ExecContext(ctx)

	return err
}
