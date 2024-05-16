package repositories

import (
	"gorm.io/gorm"
	"simple-banking-app/models"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (t *TransactionRepository) CreateTransaction(transaction models.Transaction) error {
	return t.db.Create(&transaction).Error
}

func (t *TransactionRepository) GetTransaction(filter *models.Transaction) (Transaction *models.Transaction, err error) {
	return Transaction, t.db.Model(&models.Transaction{}).Where(&filter).First(&Transaction).Error
}

func (t *TransactionRepository) CreateTransactionWithWallet(
	transactionData *models.Transaction,
	walletData *models.Wallet,
) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transactionData).Error; err != nil {
			return err
		}

		walletData.TransactionReference = transactionData.Reference

		if err := tx.Create(walletData).Error; err != nil {
			return err
		}

		return nil
	})
}

func (t *TransactionRepository) Find(filter *models.Transaction) (Transaction []*models.Transaction, err error) {
	return Transaction, t.db.Model(&models.Transaction{}).Where(&filter).Find(&Transaction).Error
}
