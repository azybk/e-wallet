package domain

import (
	"context"
	"e_wallet/backend/dto"
)

type TopUp struct {
	ID      string  `db:"id"`
	UserID  int     `db:"user_id"`
	Amount  float64 `db:"amount"`
	Status  int     `db:"status"`
	SnapURL string  `db:"snap_url"`
}

type TopUpRepository interface {
	FindById(ctx context.Context, id string) (TopUp, error)
	Insert(ctx context.Context, t *TopUp) error
	Update(ctx context.Context, t *TopUp) error
}

type TopUpService interface {
	ConfirmedTopUp(ctx context.Context, id string) error
	InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error)
}