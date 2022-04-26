package test

import (
	"jackbot/db/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var MockGame = models.Game{
	Name:         "jacken",
	Jackpot:      0,
	Numbers:      5,
	NumbersRange: 10,
	BonusNumbers: 2,
	BonusRange:   5,
	EntryFee:     5,
	Active:       true,
}

var MockUser = models.User{
	Id:          "abc123",
	Name:        "test",
	TotalAmount: 0,
}

var MockPermission = models.Permission{
	UserId:      "abc123",
	Permissions: pq.Int32Array{models.MasterAdmin},
}

func SeedGame(game *models.Game, db *gorm.DB) error {
	res := db.Create(game)
	return res.Error
}

func SeedUser(user *models.User, db *gorm.DB) error {
	res := db.Create(user)
	return res.Error
}

func SeedPermission(permission *models.Permission, db *gorm.DB) error {
	res := db.Create(permission)
	return res.Error
}
