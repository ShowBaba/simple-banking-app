package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"simple-banking-app/internal/routes"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to a simple banking app API")
}

func Routes(app *fiber.App, database *gorm.DB) {
	apiURL := "/"
	router := app.Group(apiURL)

	app.Get(apiURL, welcome)
	routes.RegisterAuthRoutes(router, database)
	routes.RegisterTransactionRoutes(router, database)
	routes.RegisterUserRoutes(router, database)
}
