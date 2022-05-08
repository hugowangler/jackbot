package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"jackbot/db"
	"jackbot/db/models"
	"jackbot/internal/bot"
	"jackbot/internal/utils"
	mathrand "math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logger := utils.NewLogger()

	logger.Info("starting jackbot")

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		logger.Fatal("missing required environment variable BOT_TOKEN")
	}

	prefix := os.Getenv("BOT_PREFIX")
	if prefix == "" {
		logger.Fatal("missing required environment variable BOT_PREFIX")
	}

	gormDb, err := db.NewConn()
	if err != nil {
		logger.With("error", err).Fatal("failed to connect to db")
	}

	var b [8]byte
	_, err = cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	mathrand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	game, err := models.GetCurrentGame(gormDb, logger)
	if err != nil {
		panic(err)
	}

	jackbot := bot.NewBot(token, prefix, game, logger, gormDb)
	err = jackbot.Start()
	if err != nil {
		logger.With("error", err).Fatal("failed to start jackbot")
	}
	defer jackbot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
		logger.Info("stopping jackbot")
		return
	}
}
