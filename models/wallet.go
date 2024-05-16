package models

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Wallet struct {
	gorm.Model `json:"-"`
	mu         sync.Mutex

	ID                   uint    `json:"-" gorm:"primarykey"`
	UserID               *uint   `json:"-" gorm:"index"`
	Balance              float64 `gorm:"not null;type:numeric(10,2);"`
	PrevBalance          float64 `json:"" gorm:"not null;type:numeric(10,2);"`
	TransactionReference string
	CreatedAt            time.Time
	UpdatedAt            time.Time `json:"-"`
}

func (w *Wallet) TopUpWalletBalance(val float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.PrevBalance = w.Balance
	w.Balance += val
}

func (w *Wallet) DeductWalletBalance(val float64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.PrevBalance = w.Balance
	w.Balance -= val

	if w.Balance < 0 {
		return errors.New("balance cannot be less than zero")
	}

	return nil
}
