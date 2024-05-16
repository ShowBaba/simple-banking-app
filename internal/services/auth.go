package services

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"simple-banking-app/internal/common"
	"simple-banking-app/internal/dtos"
	"simple-banking-app/internal/repositories"
	"simple-banking-app/internal/utils"
	"simple-banking-app/models"
)

type AuthClient interface {
	Login(c *fiber.Ctx) (*dtos.LoginResp, *common.RestErr)
}

type AuthService struct {
	userRepo repositories.UserRepository
	restErr  *common.RestErr
}

func NewAuthService(
	userRepo repositories.UserRepository,
	restErr *common.RestErr,
) AuthClient {
	return &AuthService{
		userRepo,
		restErr,
	}
}

func (a *AuthService) Login(c *fiber.Ctx) (*dtos.LoginResp, *common.RestErr) {
	var input dtos.LoginDTO

	i := c.Locals("input")
	input, ok := i.(dtos.LoginDTO)
	if !ok {
		log.Error(fmt.Errorf("cannot convert validated data to LoginDTO"))
		return nil, a.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	user, exist, err := a.userRepo.FetchOne(models.User{Email: input.Email})
	if err != nil {
		log.Error(zap.Error(err))
		return nil, a.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	if !exist {
		if err != nil {
			return nil, a.restErr.BadRequest(common.ErrUserWithEmailNotFound)
		}
	}

	passwordMatch, err := utils.PasswordMatches(input.Password, user.Password)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, a.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	if !passwordMatch {
		return nil, a.restErr.BadRequest(common.ErrInvalidPassword)
	}

	token, err := utils.GenerateToken(utils.GetConfig().JWTSecretKey, input.Email, user.ID)
	if err != nil {
		log.Error(zap.Error(err))
		return nil, a.restErr.ServerError(common.ErrSomethingWentWrong)
	}

	return &dtos.LoginResp{Token: token}, nil
}
