package handlers

import (
	"github.com/gofiber/fiber/v2"
	"simple-banking-app/internal/services"
)

type TransactionHandler struct {
	transactionSvc services.TransactionClient
}

func NewTransactionHandler(
	transactionSvc services.TransactionClient,
) *TransactionHandler {
	return &TransactionHandler{
		transactionSvc,
	}
}

func (t *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	resp, err := t.transactionSvc.CreateTransaction(c)
	if err != nil {
		return c.Status(err.StatusCode).JSON(err)
	}

	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "transaction processed successful",
		"data":    resp,
	})
}
