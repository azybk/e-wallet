package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUser(userRepository domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s userService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := s.userRepository.FindByUserName(ctx, req.Username)
	if err != nil {
		return dto.AuthRes{}, err
	}

	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GenerateRandomString(16)
	return dto.AuthRes{
		Token: token,
	}, nil
}

func (s userService) ValidateToken(ctx context.Context, token string) (userData dto.UserData, err error) {
	return
}
