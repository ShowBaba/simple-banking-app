package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/handlers"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/services"
	"simple-banking-app/internal/validators"
)

func RegisterAuthRoutes(router fiber.Router, db *gorm.DB) {
	restErr := common.NewRestErr()
	userRepo := repositories.NewUserRepository(db)
	authSvc := services.NewAuthService(*userRepo, restErr)
	validator := validators.NewAuthValidator(*userRepo, restErr)
	handler := handlers.NewAuthHandler(authSvc)

	userRouter := router.Group("auth")
	userRouter.Post("/login", validator.ValidateLogin, handler.Login)
}
