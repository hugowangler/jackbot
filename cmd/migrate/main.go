package main

import (
	"github.com/joho/godotenv"
	gormdb "jackbot/db"
	"jackbot/db/migrate"
	"jackbot/internal/utils"
)

func main() {
	_ = godotenv.Load()
	logger := utils.NewLogger()

	logger.Info("running migrations")

	db, err := gormdb.NewConn()
	if err != nil {
		logger.With("error", err).Fatal("failed to connect to db")
	}
	err = migrate.Migrate(db)
	if err != nil {
		logger.With("error", err).Fatal("failed to run migrations")
	}
	logger.Info("migration complete, db is up to date")
}
