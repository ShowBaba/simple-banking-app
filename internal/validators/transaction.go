package validators

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/dtos"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/utils"
	"simple-banking-app/models"
	"strconv"
)

type TransactionValidator struct {
	walletRepo repositories.WalletRepository
	restErr    *common.RestErr
}

func NewTransactionValidator(
	walletRepo repositories.WalletRepository,
	restErr *common.RestErr,
) *TransactionValidator {
	return &TransactionValidator{
		walletRepo,
		restErr,
	}
}

func (t *TransactionValidator) ValidateCreateTransaction(c *fiber.Ctx) error {
	var input dtos.TransactionDTO
	if err := c.BodyParser(&input); err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(t.restErr.ServerError(common.ErrFailToParseReqBody))
	}

	err := Validator.Struct(input)
	if err != nil {
		return utils.SchemaError(c, err)
	}

	if !input.IsValidTransactionType() {
		return c.Status(http.StatusBadRequest).JSON(t.restErr.ServerError(common.ErrInvalidTransactionType))
	}

	userId := c.GetRespHeader("id")
	intUserId, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Error(zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(t.restErr.ServerError(common.ErrSomethingWentWrong))
	}
	uintUserId := uint(intUserId)

	wallet, err := t.walletRepo.GetWallet(&models.Wallet{UserID: &uintUserId})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := t.walletRepo.CreateWallet(&models.Wallet{
				UserID:    &uintUserId,
				Balance:   0,
				AccountID: utils.GenerateAccountID(),
			}); err != nil {
				log.Error(zap.Error(err))
				return c.Status(http.StatusInternalServerError).JSON(t.restErr.ServerError(common.ErrSomethingWentWrong))
			}
		} else {
			log.Error(zap.Error(err))
			return c.Status(http.StatusInternalServerError).JSON(t.restErr.ServerError(common.ErrSomethingWentWrong))
		}
	}

	// validate wallet balance if transaction type is debit
	if input.Type == string(common.Debit) {

		if wallet.Balance < input.Amount {
			return c.Status(http.StatusBadRequest).JSON(t.restErr.ServerError(common.ErrInsufficientFunds))
		}
	}
	c.Locals("input", input)
	return c.Next()
}
