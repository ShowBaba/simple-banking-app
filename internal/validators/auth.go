package validators

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"net/http"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/dtos"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/utils"
	"simple-banking-app/models"
)

type AuthValidator struct {
	userRepo repositories.UserRepository
	restErr  *common.RestErr
}

func NewAuthValidator(
	userRepo repositories.UserRepository,
	restErr *common.RestErr,
) *AuthValidator {
	return &AuthValidator{
		userRepo,
		restErr,
	}
}

var Validator = validator.New()

func (a *AuthValidator) ValidateLogin(c *fiber.Ctx) error {
	var input dtos.LoginDTO
	if err := c.BodyParser(&input); err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(a.restErr.ServerError(common.ErrBadRequest))
	}

	err := Validator.Struct(input)
	if err != nil {
		return utils.SchemaError(c, err)
	}
	// validate user exist
	exist, err := a.userEmailExist(input.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(a.restErr.ServerError(common.ErrSomethingWentWrong))
	}
	if !exist {
		return c.Status(http.StatusBadRequest).JSON(a.restErr.ServerError(common.ErrUserWithEmailNotFound))
	}
	c.Locals("input", input)
	return c.Next()
}

func (a *AuthValidator) userEmailExist(email string) (bool, error) {
	_, exist, err := a.userRepo.FetchOne(models.User{Email: email})
	if err != nil {
		return false, err
	}
	return exist, nil
}
