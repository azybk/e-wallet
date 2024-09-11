package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"

	"github.com/google/uuid"
)

type topUpService struct {
	notificationService domain.NotificationService
	midtransService domain.MidtransService
	topUpRepository domain.TopUpRepository
}

func NewTopUp(notificationService domain.NotificationService, 
	midtransService domain.MidtransService,
	topUpRepository domain.TopUpRepository) domain.TopUpService {

		return &topUpService{
			notificationService: notificationService,
			midtransService: midtransService,
			topUpRepository: topUpRepository,
		}
}

func (t topUpService) InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error) {
	topUp := domain.TopUp{
		ID: uuid.NewString(),
		UserID: int(req.UserId),
		Status: 0,
		Amount: req.Amount,
	}

	err := t.midtransService.GenerateSnapURL(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	err = t.topUpRepository.Insert(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	return dto.TopUpRes{
		SnapURL: topUp.SnapURL,
	}, nil
}

func (t topUpService) ConfirmedTopUp(ctx context.Context, id string) error {
	
}