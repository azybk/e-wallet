package api

import (
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

type topUpApi struct {
	topUpService domain.TopUpService
}

func NewTopUp(app *fiber.App, authMid fiber.Handler, topUpService domain.TopUpService) {
	t := topUpApi{
		topUpService: topUpService,
	}

	app.Post("/topup/initialize", authMid, t.InitializeTopUp)
}

func (t topUpApi) InitializeTopUp(ctx *fiber.Ctx) error {
	var req dto.TopUpReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-user").(dto.UserData)
	req.UserId = user.ID

	res, err := t.topUpService.InitializeTopUp(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.ErrorType(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(res)
}
