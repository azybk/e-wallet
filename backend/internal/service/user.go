package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository domain.UserRepository
	cacheRepository domain.CacheRepository
}

func NewUser(userRepository domain.UserRepository, cacheRepository domain.CacheRepository) domain.UserService {
	return &userService{
		userRepository: userRepository,
		cacheRepository: cacheRepository,
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

	userJson, _ := json.Marshal(user)
	_ = s.cacheRepository.Set("user:" + token, userJson)

	return dto.AuthRes{
		Token: token,
	}, nil
		
}

func (s userService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := s.cacheRepository.Get("user:" + token)

	if err != nil {
		return dto.UserData{}, domain.ErrAuthFailed
	}

	var user domain.User
	_ = json.Unmarshal(data, &user)

	return dto.UserData{
		ID: user.ID,
		FullName: user.FullName,
		Phone: user.Phone,
		Username: user.Username,
	}, nil
}
