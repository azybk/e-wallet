package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"
	"encoding/json"
	"fmt"
	"time"
)

type transactionService struct {
	accountRepository domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository domain.CacheRepository
	notificationService domain.NotificationService
	// notificationRepository domain.NotificationRepository
	// hub *dto.Hub
}

func NewTransaction(accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
	notificationService domain.NotificationService,
	// notificationRepository domain.NotificationRepository,
	// hub *dto.Hub
	) domain.TransactionService {

	return &transactionService {
		accountRepository: accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository: cacheRepository,
		notificationService: notificationService,
		// notificationRepository: notificationRepository,
		// hub: hub,
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

	inquiryKey := util.GenerateRandomNumber(32)
	jsonData, _ := json.Marshal(req)

	_ = t.cacheRepository.Set(inquiryKey, jsonData)

	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

func (t transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	val, err := t.cacheRepository.Get(req.InquiryKey)
	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInq dto.TransferInquiryReq
	_ = json.Unmarshal(val, &reqInq)
	if reqInq == (dto.TransferInquiryReq{}) {
		return domain.ErrInquiryNotFound
	}

	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := t.accountRepository.FindByUserId(ctx, user.ID)
	if err != nil {
		return err
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, reqInq.AccountNumber)
	if err != nil {
		return err
	}

	debitTransaction := domain.Transaction {
		AccountId: myAccount.ID,
		SofNumber: myAccount.AccountNumber,
		DofNumber: dofAccount.AccountNumber,
		TransactionType: "D",
		Amount: reqInq.Amount,
		TransactionDatetime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &debitTransaction)
	if err != nil {
		return err
	}

	creditTransaction := domain.Transaction{
		AccountId: dofAccount.ID,
		SofNumber: myAccount.AccountNumber,
		DofNumber: dofAccount.AccountNumber,
		TransactionType: "C",
		Amount: reqInq.Amount,
		TransactionDatetime: time.Now(),
	}
	err = t.transactionRepository.Insert(ctx, &creditTransaction)
	if err != nil {
		return err
	}

	myAccount.Balance -= reqInq.Amount
	err = t.accountRepository.Update(ctx, &myAccount)
	if err != nil {
		return err
	}

	dofAccount.Balance += reqInq.Amount
	err = t.accountRepository.Update(ctx, &dofAccount)
	if err != nil {
		return err
	}

	go t.notificationAfterTransfer(myAccount, dofAccount, reqInq.Amount)
	return nil
}

func (t transactionService) notificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", amount),
	}

	t.notificationService.Insert(context.Background(), sofAccount.UserId, "TRANSFER", data)
	t.notificationService.Insert(context.Background(), dofAccount.UserId, "TRANSFER_DEST", data)

	// notificationSender := domain.Notification{
	// 	UserId: sofAccount.UserId,
	// 	Title: "Transfer berhasil",
	// 	Body: fmt.Sprintf("Transfer sejumlah %.2f berhasil", amount),
	// 	IsRead: 0,
	// 	Status: 1,
	// 	CreatedAt: time.Now(),
	// }

	// notificationReceiver := domain.Notification{
	// 	UserId: dofAccount.UserId,
	// 	Title: "Dana diterima",
	// 	Body: fmt.Sprintf("Dana diterima senilai %.2f", amount),
	// 	IsRead: 0,
	// 	Status: 1,
	// 	CreatedAt: time.Now(),
	// }
	
	// _ = t.notificationRepository.Insert(context.Background(), &notificationSender)
	// if channel, ok := t.hub.NotificationChannel[sofAccount.UserId]; ok {
	// 	channel <- dto.NotificationData{
	// 		ID: notificationSender.ID,
	// 		Title: notificationSender.Title,
	// 		Body: notificationSender.Body,
	// 		Status: notificationSender.Status,
	// 		IsRead: notificationSender.IsRead,
	// 		CreatedAt: notificationSender.CreatedAt,
	// 	}
	// }

	// _ = t.notificationRepository.Insert(context.Background(), &notificationReceiver)
	// if channel, ok := t.hub.NotificationChannel[dofAccount.UserId]; ok {
	// 	channel <- dto.NotificationData{
	// 		ID: notificationReceiver.ID,
	// 		Title: notificationReceiver.Title,
	// 		Body: notificationReceiver.Body,
	// 		Status: notificationReceiver.Status,
	// 		IsRead: notificationReceiver.IsRead,
	// 		CreatedAt: notificationReceiver.CreatedAt,
	// 	}
	// }
}