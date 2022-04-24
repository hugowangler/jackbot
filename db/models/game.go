package models

import "time"

type Game struct {
	Id           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Jackpot      int
	Numbers      int
	NumbersRange int
	BonusNumbers int
	BonusRange   int
	EntryFee     int
	Active       bool
}
