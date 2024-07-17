package api

import (
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

type authApi struct {
	userService domain.UserService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMiddleware fiber.Handler) {
	handler := authApi{
		userService: userService,
	}

	app.Post("/token/generate", handler.GenerateToken)
	app.Get("token/validate", authMiddleware, handler.ValidateToken)
}

func (a authApi) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	token, err := a.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.ErrorType(err))
	}

	return ctx.Status(200).JSON(token)
}

func (a authApi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")

	return ctx.Status(200).JSON(user)
}