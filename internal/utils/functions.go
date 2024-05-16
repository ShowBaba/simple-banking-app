package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
)

func PublishMessageToQueue(ctx context.Context, conn *amqp091.Connection, message []byte, queueName string) error {
	fmt.Println("publishing to ", queueName)
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q.Name, q.Name, queueName, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, queueName, "", false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	},
	)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// GenerateToken generates a jwt token
func GenerateToken(JWTSecretKey, email string, id uint) (signedToken string, err error) {
	claims := &AuthTokenJwtClaim{
		Email: email,
		ID:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return
	}
	return
}

func BoolPointer(b bool) *bool {
	return &b
}

func PasswordMatches(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateAuthToken(signedToken, SecretKey string) (*AuthTokenJwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&AuthTokenJwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AuthTokenJwtClaim)
	if !ok {
		return nil, err
	}
	// check the expiration date of the token
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, err
	}
	return claims, nil
}

func IsTokenValid(token TokenStruct) bool {
	currentTime := time.Now()
	duration := currentTime.Sub(token.CreatedAt)
	return duration <= 30*time.Minute
}

func GenerateReference() string {
	return fmt.Sprintf("txn-%d", time.Now().UnixNano())
}

type IError struct {
	Field string
	Tag   string
	Value string
}

func SchemaError(c *fiber.Ctx, err error) error {
	var errors []*IError
	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = fmt.Sprintf(`%s is %s`, err.Field(), err.Tag())
		errors = append(errors, &el)
	}
	return c.Status(fiber.StatusBadRequest).JSON(
		&fiber.Map{"errors": errors},
	)
}

func GetAuthUserIdFromContext(c *fiber.Ctx) (uint, error) {
	userId := c.GetRespHeader("id")
	intUserId, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(intUserId), nil
}
