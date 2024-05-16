package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/handlers"
	"simple-banking-app/internal/middleware"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/services"
)

func RegisterUserRoutes(router fiber.Router, db *gorm.DB) {
	restErr := common.NewRestErr()
	authMiddleware := middleware.NewAuthMiddleware(restErr)
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	userSvc := services.NewUserService(*userRepo, *walletRepo, *transactionRepo, restErr)
	handler := handlers.NewUserHandler(userSvc)

	userRouter := router.Group("user")
	userRouter.Get("/get-details", authMiddleware.ValidateAuthHeaderToken, handler.GetUserDetails)
}
