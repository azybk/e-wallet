package api

import (
	"e_wallet/backend/domain"

	"github.com/gofiber/fiber/v2"
)

type topUpApi struct {
	topUpService domain.TopUpService
}

func NewTopUp(app *fiber.Ctx, authMid fiber.Handler, topUpService domain.TopUpService) {
	
}