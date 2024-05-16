package handlers

import (
	"github.com/gofiber/fiber/v2"
	"simple-banking-app/internal/services"
)

type UserHandler struct {
	userSvc services.UserClient
}

func NewUserHandler(
	userSvc services.UserClient,
) *UserHandler {
	return &UserHandler{
		userSvc,
	}
}

func (u *UserHandler) GetUserDetails(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	resp, err := u.userSvc.GetUserDetails(c)
	if err != nil {
		return c.Status(err.StatusCode).JSON(err)
	}

	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "fetched user data successful",
		"data":    resp,
	})
}
