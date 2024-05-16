package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/handlers"
	"simple-banking-app/internal/middleware"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/services"
	"simple-banking-app/internal/validators"
)

func RegisterTransactionRoutes(router fiber.Router, db *gorm.DB) {
	restErr := common.NewRestErr()
	authMiddleware := middleware.NewAuthMiddleware(restErr)
	accountRepo := repositories.NewWalletRepository(db)
	validator := validators.NewTransactionValidator(*accountRepo, restErr)
	transactionRepo := repositories.NewTransactionRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionSvc := services.NewTransactionService(*transactionRepo, *walletRepo, restErr)
	handler := handlers.NewTransactionHandler(transactionSvc)

	userRouter := router.Group("transaction")
	userRouter.Post("/create-transaction", authMiddleware.ValidateAuthHeaderToken, validator.ValidateCreateTransaction, handler.CreateTransaction)
}
