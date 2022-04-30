package test

import (
	"jackbot/db/models"

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
	AccountantId: MockUser.Id,
}

func SeedGame(game *models.Game, db *gorm.DB) error {
	res := db.Create(game)
	return res.Error
}
