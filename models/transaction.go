package models

import (
	"gorm.io/gorm"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/utils"
	"time"
)

type Transaction struct {
	ID        uint                   `gorm:"primaryKey"`
	Type      common.TransactionType `gorm:"type:varchar(6);not null"`
	Amount    float64                `gorm:"not null;type:numeric(10,2);"`
	Narration string
	Status    common.TransactionStatus `gorm:"type:varchar(10);not null"`
	Reference string                   `gorm:"uniqueIndex;not null"`
	UserID    *uint                    `json:"-" gorm:"index"`
	AccountID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.Reference = utils.GenerateReference()
	return
}
