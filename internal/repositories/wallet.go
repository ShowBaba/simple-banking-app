package repositories

import (
	"gorm.io/gorm"
	"simple-banking-app/models"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db}
}

func (w *WalletRepository) GetWallet(filter *models.Wallet) (Wallet *models.Wallet, err error) {

	return Wallet, w.db.Model(&models.Wallet{}).Where(&filter).
		Last(&Wallet).Error
}

func (w *WalletRepository) CreateWallet(wallet *models.Wallet) error {
	return w.db.Create(&wallet).Error
}
