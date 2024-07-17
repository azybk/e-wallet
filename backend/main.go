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
	userService := service.NewUser(userRepository, cacheConnection)

	authMiddleware := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMiddleware)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}