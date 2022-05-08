package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Row struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RaffleId     uint64 `gorm:"primaryKey"`
	Raffle       Raffle
	Numbers      pq.Int32Array `gorm:"type:integer[]"`
	BonusNumbers pq.Int32Array `gorm:"type:integer[]"`
	UserId       string        `gorm:"primaryKey"`
	User         User
	Paid         bool
}

func (r *Row) NumbersToString() string {
	var res string
	for _, b := range r.Numbers {
		res += fmt.Sprintf("%d ", b)
	}
	return strings.TrimSpace(res)
}

func (r *Row) BonusNumbersToString() string {
	var res string
	for _, b := range r.BonusNumbers {
		res += fmt.Sprintf("%d ", b)
	}
	return strings.TrimSpace(res)
}
