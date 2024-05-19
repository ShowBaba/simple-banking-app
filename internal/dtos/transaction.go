package dtos

import "simple-banking-app/internal/common"

type TransactionDTO struct {
	Type      string  `json:"type" validate:"required"`
	Amount    float64 `json:"amount" validate:"required"`
	Narration string  `json:"narration"`
}

var validTransactionType = map[common.TransactionType]bool{
	common.Debit:  true,
	common.Credit: true,
}

func (t *TransactionDTO) IsValidTransactionType() bool {
	return validTransactionType[common.TransactionType(t.Type)]
}
