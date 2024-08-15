package main

import (
	"e_wallet/backend/internal/api"
	"e_wallet/backend/internal/component"
	"e_wallet/backend/internal/config"
	"e_wallet/backend/internal/middleware"
	"e_wallet/backend/internal/repository"
	"e_wallet/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	cacheConnection := component.GetCacheConnection()

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection)

	authMiddleware := middleware.Authenticate(userService)
	
	app := fiber.New()
	api.NewAuth(app, userService, authMiddleware)
	api.NewTransfer(app, authMiddleware, transactionService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
