package util

import (
	"e_wallet/backend/domain"
	"errors"
)

func ErrorType(err error) int {
	switch {
		case errors.Is(err, domain.ErrAuthFailed):
			return 401
		case errors.Is(err, domain.ErrUsernameTaken):
			return 400
		case errors.Is(err, domain.ErrOtpInvalid):
			return 400
		default:
			return 500
	}
}