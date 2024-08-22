package api

import (
	"context"
	"e_wallet/backend/domain"
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

type notificationApi struct {
	notificationService domain.NotificationService
}

func NewNotification(app *fiber.App, authMid fiber.Handler, notificationService domain.NotificationService) {
	h := notificationApi{
		notificationService: notificationService,
	}

	app.Get("/notifications", authMid, h.GetUserNotification)
}

func (n notificationApi) GetUserNotification(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15 * time.Second)
	defer cancel()

	user := ctx.Locals("x-user").(dto.UserData)

	notification, err := n.notificationService.FindByUser(c, user.ID)
	if err != nil {
		return ctx.Status(util.ErrorType(err)).JSON(dto.Response{Message: err.Error()})
	}

	return ctx.Status(200).JSON(notification)
}