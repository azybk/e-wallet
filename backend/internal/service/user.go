package service

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository  domain.UserRepository
	cacheRepository domain.CacheRepository
	emailService    domain.EmailService
}

func NewUser(userRepository domain.UserRepository, cacheRepository domain.CacheRepository, emailService domain.EmailService) domain.UserService {
	return &userService{
		userRepository:  userRepository,
		cacheRepository: cacheRepository,
		emailService:    emailService,
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

	if !user.EmailVerifiedAtDB.Valid {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GenerateRandomString(16)

	userJson, _ := json.Marshal(user)
	_ = s.cacheRepository.Set("user:"+token, userJson)

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
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}

func (s userService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	exist, err := s.userRepository.FindByUserName(ctx, req.Username)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	if exist != (domain.User{}) {
		return dto.UserRegisterRes{}, domain.ErrUsernameTaken
	}

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	err = s.userRepository.Insert(ctx, &user)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	otpCode := util.GenerateRandomNumber(4)
	referenceId := util.GenerateRandomString(16)

	s.emailService.Send(req.Email, "OTP Code", "OTP anda "+otpCode)

	_ = s.cacheRepository.Set("otp:"+referenceId, []byte(otpCode))
	_ = s.cacheRepository.Set("user-ref:"+referenceId, []byte(user.Username))

	return dto.UserRegisterRes{
		ReferenceID: referenceId,
	}, nil
}

func (s userService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	val, err := s.cacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}

	otp := string(val)
	if otp != req.OTP {
		return domain.ErrOtpInvalid
	}

	val, err = s.cacheRepository.Get("user-ref:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}

	user, err := s.userRepository.FindByUserName(ctx, string(val))
	if err != nil {
		return err
	}

	user.EmailVerifiedAt = time.Now()
	_ = s.userRepository.Update(ctx, &user)

	return nil
}
