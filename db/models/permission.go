package models

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Permission struct {
	UserId      string        `gorm:"primaryKey"`
	Permissions pq.Int32Array `gorm:"type:integer[]"`
}

const (
	MasterAdmin = iota + 1
)

func HasPermissions(userId string, permissions []int, db *gorm.DB) (bool, error) {
	var permissionModel Permission
	res := db.Where("user_id = ? AND permissions <@ ?", userId, pq.Array(permissions)).First(&permissionModel)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, res.Error
	}

	return true, nil
}
