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
	user := ctx.Value("x-user").(dto.UserData)

	myAccount, err := t.accountRepository.FindByUserId(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	if myAccount.Balance < req.Amount {
		return dto.TransferInquiryRes{}, domain.ErrInsuficientBlanace
	}

}

func (t transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {

}