package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/internal/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransService struct {
	client snap.Client
}

func NewMidtransService(cnf *config.Config) domain.MidtransService {
	var client snap.Client
	
	envi := midtrans.Sandbox
	if cnf.Midtrans.IsProd {
		envi = midtrans.Production
	}

	client.New(cnf.Midtrans.Key, envi)

	return &midtransService{
		client: snap,
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

}