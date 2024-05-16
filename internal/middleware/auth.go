package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/utils"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	restErr *common.RestErr
}

func NewAuthMiddleware(
	restErr *common.RestErr,
) *AuthMiddleware {
	return &AuthMiddleware{
		restErr,
	}
}

func (a *AuthMiddleware) ValidateAuthHeaderToken(c *fiber.Ctx) error {
	tokenInHeader := c.Get("Authorization")
	if tokenInHeader == "" {
		return c.Status(http.StatusBadRequest).JSON(a.restErr.ServerError(common.ErrMissingAuthTokenInHeader))
	}
	token := strings.Split(tokenInHeader, " ")[1]
	if token == "" {
		return c.Status(http.StatusBadRequest).JSON(a.restErr.ServerError(common.ErrMissingAuthTokenInHeader))
	}
	claim, err := utils.ValidateAuthToken(token, utils.GetConfig().JWTSecretKey)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(a.restErr.ServerError(common.ErrInvalidAuthToken))
	}

	// TODO: fetch user and confirm if email has been validated
	c.Set("email", claim.Email)
	c.Set("id", strconv.FormatUint(uint64(claim.ID), 10))
	return c.Next()
}
