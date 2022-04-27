package migrate

import (
	"jackbot/db/models"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	m := gormigrate.New(
		db, gormigrate.DefaultOptions, []*gormigrate.Migration{
			{
				ID: "202204242311",
				Migrate: func(tx *gorm.DB) error {
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
					}
					type Raffle struct {
						Id           uint64 `gorm:"primaryKey"`
						CreatedAt    time.Time
						UpdatedAt    time.Time
						GameId       uint64
						Game         models.Game
						WinningRowId uint64
						Date         time.Time
					}
					type Row struct {
						Id           uint64 `gorm:"primaryKey"`
						CreatedAt    time.Time
						UpdatedAt    time.Time
						RaffleId     uint64
						Raffle       Raffle
						Numbers      pq.Int32Array `gorm:"type:integer[]"`
						BonusNumbers pq.Int32Array `gorm:"type:integer[]"`
						UserId       string
						User         models.User
						Paid         bool
					}
					type User struct {
						Id          string `gorm:"primaryKey"`
						CreatedAt   time.Time
						UpdatedAt   time.Time
						Name        string
						TotalAmount int
					}
					type Permission struct {
						UserId      string        `gorm:"primaryKey"`
						Permissions pq.Int32Array `gorm:"type:integer[]"`
					}

					return tx.AutoMigrate(&Game{}, &Raffle{}, &Row{}, &User{}, &Permission{})
				},
				Rollback: func(tx *gorm.DB) error {
					if err := tx.Migrator().DropTable("games"); err != nil {
						return err
					}

					if err := tx.Migrator().DropTable("raffles"); err != nil {
						return err
					}

					if err := tx.Migrator().DropTable("rows"); err != nil {
						return err
					}

					if err := tx.Migrator().DropTable("users"); err != nil {
						return err
					}

					if err := tx.Migrator().DropTable("permissions"); err != nil {
						return err
					}

					return nil
				},
			},
			{
				ID: "202204242353",
				Migrate: func(tx *gorm.DB) error {
					admins := []models.User{
						{
							Id:          "178216786666323968",
							Name:        "Hugo Wangler",
							TotalAmount: 0,
						},
						{
							Id:          "178632146762596352",
							Name:        "Måns Falk",
							TotalAmount: 0,
						},
					}

					if err := tx.Create(&admins).Error; err != nil {
						return err
					}

					permissions := []models.Permission{
						{
							UserId:      "178216786666323968",
							Permissions: pq.Int32Array([]int32{models.MasterAdmin}),
						},
						{
							UserId:      "178632146762596352",
							Permissions: pq.Int32Array([]int32{models.MasterAdmin}),
						},
					}

					if err := tx.Create(&permissions).Error; err != nil {
						return err
					}

					return tx.AutoMigrate(&admins, &permissions)
				},
				Rollback: func(tx *gorm.DB) error {
					var err error

					if err = tx.Delete(&models.User{}, []string{"178216786666323968", "178632146762596352"}).Error; err != nil {
						return err
					}

					if err = tx.Delete(&models.Permission{}, []string{"178216786666323968", "178632146762596352"}).Error; err != nil {
						return err
					}

					return nil
				},
			},
		},
	)
	err := m.Migrate()
	return err
}
