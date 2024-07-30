package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
)

type transactionService struct {
	accountRepository domain.AccountRepository
	transacionRepository domain.TransactionRepository
	cacheRepository domain.CacheRepository
}

func NewTransaction(accountRepository domain.AccountRepository,
	transacionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository) domain.TransactionService {

	return &transactionService {
		accountRepository: accountRepository,
		transacionRepository: transacionRepository,
		cacheRepository: cacheRepository,
	}
}

func (t transactionService) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {

}

func (t transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {

}