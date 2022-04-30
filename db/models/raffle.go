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
	Date         *time.Time
}

func CreateRaffle(raffle *Raffle, db *gorm.DB) error {
	existingRaffle := []Raffle{}
	if err := db.Where("date IS NULL AND game_id = ?", raffle.GameId).Preload("Game").Find(&existingRaffle).Error; err != nil {
		return err
	}

	if len(existingRaffle) > 0 {
		return &PreviousRaffleNotCompletedError{name: existingRaffle[0].Game.Name}
	}

	return db.Create(raffle).Error
}
