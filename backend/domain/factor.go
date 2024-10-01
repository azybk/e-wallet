package domain

import (
	"context"
	"e_wallet/backend/dto"
)

type Factor struct {
	ID     int64  `db:"id"`
	UserID int64  `db:"user_id"`
	pin    string `db:"pin"`
}

type FactorRepository interface {
	FindByUser(ctx context.Context, id int64) (Factor, error)
}

type FactorService interface {
	ValidatePIN(ctx context.Context, req dto.ValidatePinReq) error
}