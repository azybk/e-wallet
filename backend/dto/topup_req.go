package dto

type TopUpReq struct {
	Amount float64 `json:"amount"`
	UserId int64   `json:"-"`
}