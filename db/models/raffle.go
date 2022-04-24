package models

import "time"

type Raffle struct {
	Id           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	GameId       uint64
	Game         Game
	WinningRowId uint64
	Date         time.Time
}
