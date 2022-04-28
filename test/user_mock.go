package test

import (
	"jackbot/db/models"

	"gorm.io/gorm"
)

var MockUser = models.User{
	Id:          "abc123",
	Name:        "test",
	TotalAmount: 0,
}

func SeedUser(user *models.User, db *gorm.DB) error {
	res := db.Create(user)
	return res.Error
}
