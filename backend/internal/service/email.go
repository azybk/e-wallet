package service

import (
	"e_wallet/backend/domain"
	"e_wallet/backend/internal/config"
)

type emailService struct {
	cnf *config.Config
}

func NewEmail(cnf *config.Config) domain.EmailService {
	return &emailService{cnf: cnf}
}

func (e emailService) Send(to, subject, body string) error {
	
}