package test

import (
	"jackbot/db/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var MockPermission = models.Permission{
	UserId:      "abc123",
	Permissions: pq.Int32Array{models.MasterAdmin},
}

func SeedPermission(permission *models.Permission, db *gorm.DB) error {
	res := db.Create(permission)
	return res.Error
}
