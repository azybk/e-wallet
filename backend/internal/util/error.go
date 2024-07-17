package util

import (
	"e_wallet/backend/domain"
	"errors"
)

func ErrorType(err error) int {
	switch {
		case errors.Is(err, domain.ErrAuthFailed):
			return 401
		default:
			return 500
	}
}