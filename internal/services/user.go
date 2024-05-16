package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/utils"
	"simple-banking-app/models"
)

type UserClient interface {
	GetUserDetails(c *fiber.Ctx) (*models.UserWithDetails, *common.RestErr)
}

type UserService struct {
	userRepo        repositories.UserRepository
	walletRepo      repositories.WalletRepository
	transactionRepo repositories.TransactionRepository
	restErr         *common.RestErr
}

func NewUserService(
	userRepo repositories.UserRepository,
	walletRepo repositories.WalletRepository,
	transactionRepo repositories.TransactionRepository,
	restErr *common.RestErr,
) UserClient {
	return &UserService{
		userRepo,
		walletRepo,
		transactionRepo,
		restErr,
	}
}

func (u *UserService) GetUserDetails(c *fiber.Ctx) (*models.UserWithDetails, *common.RestErr) {
	userId, err := utils.GetAuthUserIdFromContext(c)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, u.restErr.ServerError(common.ErrSomethingWentWrong)
	}
	user, exist, err := u.userRepo.FetchOne(models.User{ID: userId})
	if err != nil {
		log.Error(zap.Error(err))
		return nil, u.restErr.ServerError(common.ErrSomethingWentWrong)
	}
	if !exist {
		return nil, u.restErr.ServerError(common.ErrUserWithEmailNotFound)
	}

	walletData, err := u.walletRepo.GetWallet(&models.Wallet{UserID: &userId})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := u.walletRepo.CreateWallet(&models.Wallet{
				UserID:  &userId,
				Balance: 0,
			}); err != nil {
				log.Error(zap.Error(err))
				return nil, u.restErr.ServerError(common.ErrSomethingWentWrong)
			}
		} else {
			log.Error(zap.Error(err))
			return nil, u.restErr.ServerError(common.ErrSomethingWentWrong)
		}
	}

	transactions, err := u.transactionRepo.Find(&models.Transaction{UserID: &userId})
	if err != nil {
		log.Error(zap.Error(err))
		return nil, u.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	details := models.UserWithDetails{
		User:         user,
		Wallet:       walletData,
		Transactions: transactions,
	}
	return &details, nil
}
