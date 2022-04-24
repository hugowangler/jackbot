package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"time"
)

func Migrate(db *gorm.DB) error {
	m := gormigrate.New(
		db, gormigrate.DefaultOptions, []*gormigrate.Migration{
			{
				ID: "202204241140",
				Migrate: func(tx *gorm.DB) error {
					type Lol struct {
						ID        string `gorm:"primaryKey"`
						CreatedAt time.Time
						UpdatedAt time.Time
					}
					return tx.AutoMigrate(&Lol{})
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable("lols")
				},
			},
		},
	)
	err := m.Migrate()
	return err
}
