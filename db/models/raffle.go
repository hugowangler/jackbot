package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PreviousRaffleNotCompletedError struct {
	name string
}

func (e *PreviousRaffleNotCompletedError) Error() string {
	return fmt.Sprintf("a raffle in currently ongoing for the game: %s", e.name)
}

type Raffle struct {
	Id           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	GameId       uint64
	Game         Game
	WinningRowId uint64
	Date         time.Time
}

func CreateRaffle(raffle *Raffle, db *gorm.DB) error {
	existingRaffle := &Raffle{}
	if err := db.Where("date = ? AND game_id = ?", nil, raffle.GameId).First(existingRaffle).Error; err != nil {
		return err
	}

	if existingRaffle != nil {
		return &PreviousRaffleNotCompletedError{name: existingRaffle.Game.Name}
	}

	return db.Create(raffle).Error
}
