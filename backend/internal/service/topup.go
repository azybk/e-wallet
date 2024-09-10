package service

import (
	"context"
	"e_wallet/backend/domain"
)

type topUpService struct {
	notificationService domain.NotificationService
	topUpRepository domain.TopUpRepository
}

func NewTopUp(notificationService domain.NotificationService, 
	topUpRepository domain.TopUpRepository) domain.TopUpService {

		return &topUpService{
			notificationService: notificationService,
			topUpRepository: topUpRepository,
		}
}

func (t topUpService) ConfirmedTopUp(ctx context.Context, id string) error {
	
}