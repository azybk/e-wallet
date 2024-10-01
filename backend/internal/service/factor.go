package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
)

type factorService struct {
	factorRepository domain.FactorRepository
}

func NewFactor(factorRepository domain.FactorRepository) domain.FactorService {
	return &factorService{
		factorRepository: factorRepository,
	}
}

func (f factorService) ValidatePIN(ctx context.Context, req dto.ValidatePinReq) error {
	factor, err := f.factorRepository.FindByUser(ctx, req.UserID)
	if err != nil {
		return err
	}

	if factor == (domain.Factor{}) {
		return domain.ErrPinInvalid
	}
}