package models

import (
	"time"
)

type User struct {
	ID             uint `gorm:"primaryKey"`
	Email          string
	FirstName      string
	LastName       string
	Password       string
	PhoneNumber    string
	Username       string
	ProfilePicture string
	IsVerified     *bool     `gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserWithDetails struct {
	User         *User
	Wallet       *Wallet
	Transactions []*Transaction
}
