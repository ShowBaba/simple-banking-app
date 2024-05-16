package db

import (
	"gorm.io/gorm"
	"simple-banking-app/models"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
	)
}
