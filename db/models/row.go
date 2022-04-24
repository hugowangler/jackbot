package models

import (
	"time"

	"github.com/lib/pq"
)

type Row struct {
	Id           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RaffleId     uint64
	Raffle       Raffle
	Numbers      pq.Int32Array `gorm:"type:integer[]"`
	BonusNumbers pq.Int32Array `gorm:"type:integer[]"`
	UserId       string
	User         User
	Paid         bool
}
