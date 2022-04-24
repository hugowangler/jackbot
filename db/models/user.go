package models

import "time"

type User struct {
	Id          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	TotalAmount int
}
