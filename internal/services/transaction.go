package services

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/dtos"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/utils"
	"simple-banking-app/models"
)

type TransactionClient interface {
	CreateTransaction(c *fiber.Ctx) (*dtos.TransactionDTO, *common.RestErr)
}

type TransactionService struct {
	transactionRepo repositories.TransactionRepository
	walletRepo      repositories.WalletRepository
	restErr         *common.RestErr
}

func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	walletRepo repositories.WalletRepository,
	restErr *common.RestErr,
) TransactionClient {
	return &TransactionService{
		transactionRepo,
		walletRepo,
		restErr,
	}
}

func (t *TransactionService) CreateTransaction(c *fiber.Ctx) (*dtos.TransactionDTO, *common.RestErr) {
	var input dtos.TransactionDTO

	i := c.Locals("input")
	input, ok := i.(dtos.TransactionDTO)
	if !ok {
		log.Error(fmt.Errorf("cannot convert validated data to LoginDTO"))
		return nil, t.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	userId, err := utils.GetAuthUserIdFromContext(c)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, t.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	wallet, err := t.walletRepo.GetWallet(&models.Wallet{UserID: &userId})
	if err != nil {
		log.Error(zap.Error(err))
		return nil, t.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	transaction := models.Transaction{
		Status:    common.Pending,
		Type:      common.TransactionType(input.Type),
		UserID:    &userId,
		Amount:    input.Amount,
		AccountID: "default_account_id",
	}

	responseFromAPI, err := t.sendToThirdPartyProvider(&transaction)
	if err != nil {
		transaction.Status = common.Failed
		log.Info(fmt.Sprintf("error from third party api %v", err))

		// create the transaction record
		if err := t.transactionRepo.CreateTransaction(transaction); err != nil {
			log.Error(zap.Error(err))
			return nil, t.restErr.ServerError(common.ErrSomethingWentWrong)
		}

	} else {
		if responseFromAPI.Status == string(common.Failed) {
			transaction.Status = common.Failed
		} else {

			// if type debit and status complete, deduct
			if input.Type == string(common.Debit) && responseFromAPI.Status == string(common.Completed) {
				err := wallet.DeductWalletBalance(input.Amount)
				if err != nil {
					log.Error(zap.Error(err))
					return nil, t.restErr.ServerError(common.ErrSomethingWentWrong)
				}
			} else if input.Type == string(common.Credit) && responseFromAPI.Status == string(common.Completed) {
				wallet.TopUpWalletBalance(input.Amount)
			}

			transaction.Status = common.TransactionStatus(responseFromAPI.Status)

			// create the transaction along with a new wallet record
			newWallet := models.Wallet{
				UserID:      &userId,
				PrevBalance: wallet.PrevBalance,
				Balance:     wallet.Balance,
			}

			err = t.transactionRepo.CreateTransactionWithWallet(&transaction, &newWallet)
			if err != nil {
				log.Error(fmt.Errorf("error creating transaction record %v", err))
				return nil, t.restErr.ServerError(common.ErrSomethingWentWrong)
			}
		}
	}

	return &input, nil
}

func (t *TransactionService) sendToThirdPartyProvider(transaction *models.Transaction) (*ThirdPartyTransactionResponse, error) {
	// req := ThirdPartyTransactionRequest{
	// 	AccountID: transaction.AccountID,
	// 	Reference: transaction.Reference,
	// 	Amount:    transaction.Amount,
	// }
	// data, err := json.Marshal(req)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// resp, err := http.Post(fmt.Sprintf("%s/third-party/payments", utils.GetConfig().ThirdPartyTnxServiceBaseURL), "application/json", bytes.NewBuffer(data))
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, errors.New("failed to create transaction with third-party provider")
	// }
	//
	// var res ThirdPartyTransactionResponse
	// if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
	// 	return nil, err
	// }
	//
	// if res.Reference != transaction.Reference {
	// 	return nil, errors.New("transaction reference mismatch")
	// }

	return &ThirdPartyTransactionResponse{
		AccountID: transaction.AccountID,
		Reference: transaction.Reference,
		Amount:    transaction.Amount,
		Status:    string(common.Completed),
	}, nil
}

type ThirdPartyTransactionRequest struct {
	AccountID string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

type ThirdPartyTransactionResponse struct {
	AccountID string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
	Status    string
}
