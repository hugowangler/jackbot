package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
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

func InitializeDevGame(db *gorm.DB) (Game, error) {
	game := &Game{
		Name:         "DevJackbot",
		Numbers:      5,
		NumbersRange: 50,
		BonusNumbers: 2,
		BonusRange:   12,
		EntryFee:     5,
		Active:       true,
		AccountantId: "178632146762596352",
	}

	error := db.Create(game).Preload("Accountant").Error
	return *game, error
}

func GetCurrentGame(db *gorm.DB, logger *zap.SugaredLogger) (*Game, error) {
	var games []Game
	err := db.Where("active", true).Preload("Accountant").Find(&games).Error
	if err != nil {
		logger.Fatal(err)
		return nil, fmt.Errorf("failed to get games from db")
	}

	if len(games) == 0 {
		env := os.Getenv("ENVIRONMENT")
		if strings.Contains("dev", strings.ToLower(env)) {
			game, err := InitializeDevGame(db)
			if err == nil {
				games = append(games, game)
			}
		}
	}

	if len(games) == 0 {
		return nil, fmt.Errorf("no games found")
	}

	return &games[0], err
}
