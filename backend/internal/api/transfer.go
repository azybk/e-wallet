package api

import (
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

type transferApi struct {
	transactionService domain.TransactionService
	factorService domain.FactorService
}

func NewTransfer(app *fiber.App, authMid fiber.Handler, transactionService domain.TransactionService, factorService domain.FactorService) {
	h := transferApi{
		transactionService: transactionService,
		factorService: factorService,
	}

	app.Post("/transfer/inquiry", authMid, h.TransferInquiry)
	app.Post("/transfer/execute", authMid, h.TransferExecute)
}

func (t transferApi) TransferInquiry(ctx *fiber.Ctx) error {
	var req dto.TransferInquiryReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	inquiry, err := t.transactionService.TransferInquiry(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.ErrorType(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(inquiry)
}

func (t transferApi) TransferExecute(ctx *fiber.Ctx) error {
	var req dto.TransferExecuteReq

	if err:= ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: "Bad Request",
		})
	}

	user := ctx.Locals("x-user").(dto.UserData)

	if err := t.factorService.ValidatePIN(ctx.Context(), dto.ValidatePinReq{
		PIN: req.PIN,
		UserID: user.ID,
	}); err != nil {
		return ctx.Status(util.ErrorType(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	err := t.transactionService.TransferExecute(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.ErrorType(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(dto.Response{
		Message: "Successfully",
	})
}