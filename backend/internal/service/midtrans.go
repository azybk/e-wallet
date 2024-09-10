package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/internal/config"
	"errors"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransService struct {
	client snap.Client
	midtransConfig config.Midtrans
	topUpService domain.TopUpService
}

func NewMidtransService(cnf *config.Config, topUpService domain.TopUpService) domain.MidtransService {
	var client snap.Client
	
	envi := midtrans.Sandbox
	if cnf.Midtrans.IsProd {
		envi = midtrans.Production
	}

	client.New(cnf.Midtrans.Key, envi)

	return &midtransService{
		client: client,
		midtransConfig: cnf.Midtrans,
		topUpService: topUpService,
	}
}

func (m midtransService) GenerateSnapURL(ctx context.Context, t *domain.TopUp) error {
	req := & snap.Request {
		TransactionDetails: midtrans.TransactionDetails{
		OrderID:  t.ID,
		GrossAmt: int64(t.Amount),
		}, 
	}

	snapResp, err := m.client.CreateTransaction(req)
	if err != nil {
		return err
	}

	t.SnapURL = snapResp.RedirectURL
	return nil
}

func (m midtransService) VerifyPayment(ctx context.Context, data map[string] interface{}) error {
	var client coreapi.Client
	
	envi := midtrans.Sandbox
	if m.midtransConfig.IsProd {
		envi = midtrans.Production
	}

	client.New(m.midtransConfig.Key, envi)

	// 3. Get order-id from payload
	orderId, exists := data["order_id"].(string)
	if !exists {
		// do something when key `order_id` not found
		return errors.New("invalid payload")
	}

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(orderId)
	if e != nil {
		return e

	} else {

		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					m.topUpService.ConfirmedTopUp(ctx, orderId)
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				m.topUpService.ConfirmedTopUp(ctx, orderId)
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return nil
}