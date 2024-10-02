package main

import (
	"e_wallet/backend/dto"
	"e_wallet/backend/internal/api"
	"e_wallet/backend/internal/component"
	"e_wallet/backend/internal/config"
	"e_wallet/backend/internal/middleware"
	"e_wallet/backend/internal/repository"
	"e_wallet/backend/internal/service"
	"e_wallet/backend/internal/sse"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	// cacheConnection := component.GetCacheConnection()	with bigCache
	cacheConection := repository.NewRedisClient(cnf)

	hub := &dto.Hub{
		NotificationChannel: map[int64]chan dto.NotificationData{},
	}

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)
	templateRepository := repository.NewTemplate(dbConnection)
	topupRepository := repository.NewTopUp(dbConnection)
	factorRepository := repository.NewFactor(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConection, emailService)
	notificationService := service.NewNotification(notificationRepository, templateRepository, hub)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConection, notificationService)
	// transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConection, notificationRepository, hub)
	midtransService := service.NewMidtransService(cnf)
	topupService := service.NewTopUp(notificationService, midtransService, topupRepository, accountRepository, transactionRepository)
	factorService := service.NewFactor(factorRepository)

	authMiddleware := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMiddleware)
	api.NewTransfer(app, authMiddleware, transactionService, factorService)
	api.NewNotification(app, authMiddleware, notificationService)
	api.NewTopUp(app, authMiddleware, topupService)
	api.NewMidtrans(app, midtransService, topupService)

	sse.NewNotification(app, authMiddleware, hub)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
