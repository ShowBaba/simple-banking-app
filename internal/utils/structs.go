package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type AuthTokenJwtClaim struct {
	Email string
	ID    uint
	jwt.StandardClaims
}

type TokenStruct struct {
	UserID    uint
	Token     int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
