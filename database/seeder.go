package db

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"simple-banking-app/models"
	"time"
)

func StartSeeder(db *gorm.DB) error {
	return seedUser(db)
}

func seedUser(db *gorm.DB) error {
	users := []models.User{
		{
			ID:             1,
			Email:          "sam@mail.com",
			FirstName:      "Sam",
			LastName:       "Show",
			Password:       "password1",
			PhoneNumber:    "1234567890",
			Username:       "sam",
			ProfilePicture: "profile1.jpg",
			IsVerified:     boolPtr(true),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             2,
			Email:          "show@mail.com",
			FirstName:      "Show",
			LastName:       "Sam",
			Password:       "password2",
			PhoneNumber:    "1234567890",
			Username:       "show",
			ProfilePicture: "profile2.jpg",
			IsVerified:     boolPtr(true),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for i, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)

		var u models.User
		err = db.Model(&models.User{}).Where(&models.User{ID: user.ID}).First(&u).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&users[i]).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if err := db.
				Where(&models.User{ID: user.ID}).
				Save(&user).
				Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
