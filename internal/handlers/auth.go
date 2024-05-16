package handlers

import (
	"github.com/gofiber/fiber/v2"
	"simple-banking-app/internal/services"
)

type AuthHandler struct {
	authSvc services.AuthClient
}

func NewAuthHandler(authSvc services.AuthClient) *AuthHandler {
	return &AuthHandler{
		authSvc,
	}
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	resp, err := a.authSvc.Login(c)
	if err != nil {
		return c.Status(err.StatusCode).JSON(err)
	}

	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "login successful",
		"data":    resp,
	})
}
