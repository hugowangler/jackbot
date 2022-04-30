package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserAlreadyExistsError struct {
	userId string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with id %s already exists", e.userId)
}

type User struct {
	Id          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Mobile      string
	TotalAmount int
}

func GetUser(userId string, db *gorm.DB) (*User, error) {
	var existingUser *User
	if err := db.Where("Id = ?", userId).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return existingUser, nil
}

func CreateUser(user User, db *gorm.DB) error {
	var existingUser *User
	if err := db.Where("Id = ?", user.Id).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(&user).Error
		}

		return err
	}

	return &UserAlreadyExistsError{userId: user.Id}
}
