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
				ID: "202204241140",
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

					return tx.AutoMigrate(&Game{}, &Raffle{}, &Row{}, &User{})
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

					return nil
				},
			},
		},
	)
	err := m.Migrate()
	return err
}
