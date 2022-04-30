package models

import (
	"errors"
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

type NoActiveRaffleError struct {
	userId string
}

func (e *NoActiveRaffleError) Error() string {
	return fmt.Sprintf("there currently no active raffle")
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

func GetActiveRaffle(gameId uint64, db *gorm.DB) (*Raffle, error) {
	var existingRaffle *Raffle
	if err := db.Where("date IS NULL AND game_id = ?", gameId).First(&existingRaffle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &NoActiveRaffleError{}
		} else {
			return nil, err
		}
	}

	return existingRaffle, nil
}
