package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GameAlreadyExistsError struct {
	name string
}

func (e *GameAlreadyExistsError) Error() string {
	return fmt.Sprintf("game with name %s already exists", e.name)
}

type Game struct {
	Id           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Jackpot      int
	Numbers      int
	NumbersRange int
	BonusNumbers int
	BonusRange   int
	EntryFee     int
	Active       bool
	AccountantId string
	Accountant   User
}

func CreateGame(game *Game, db *gorm.DB) error {
	existingGame := &[]Game{}
	if err := db.Where("name = ?", game.Name).Find(existingGame).Error; err != nil {
		return err
	}

	if len(*existingGame) > 0 {
		return &GameAlreadyExistsError{name: game.Name}
	}

	return db.Create(game).Error
}

func GetAccountant(game *Game, db *gorm.DB) (*User, error) {
	var accountant *User
	if err := db.Where("Id = ?", game.AccountantId).First(&accountant).Error; err != nil {
		return nil, err
	}

	return accountant, nil
}
