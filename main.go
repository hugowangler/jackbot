package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"github.com/joho/godotenv"
	"jackbot/internal/bot"
	"jackbot/internal/game"
	"jackbot/utils"
	mathrand "math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	envBalls := os.Getenv("BALLS")
	if prefix == "" {
		logger.Fatal("missing required environment variable BALLS")
	}
	balls, err := strconv.Atoi(envBalls)
	if err != nil {
		logger.Fatal("BALLS environment variable was not a number")
	}

	envBonus := os.Getenv("BONUS")
	if prefix == "" {
		logger.Fatal("missing required environment variable BONUS")
	}
	bonus, err := strconv.Atoi(envBonus)
	if err != nil {
		logger.Fatal("BONUS environment variable was not a number")
	}

	envBallsRange := os.Getenv("BALLS_RANGE")
	if prefix == "" {
		logger.Fatal("missing required environment variable BALLS_RANGE")
	}
	ballsRange, err := strconv.Atoi(envBallsRange)
	if err != nil {
		logger.Fatal("BALLS_RANGE environment variable was not a number")
	}

	envBonusRange := os.Getenv("BONUS_RANGE")
	if prefix == "" {
		logger.Fatal("missing required environment variable BONUS_RANGE")
	}
	bonusRange, err := strconv.Atoi(envBonusRange)
	if err != nil {
		logger.Fatal("BONUS_RANGE environment variable was not a number")
	}

	var b [8]byte
	_, err = cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	mathrand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	g := game.Game{
		Balls:      balls,
		BallsRange: ballsRange,
		Bonus:      bonus,
		BonusRange: bonusRange,
	}

	jackbot := bot.NewBot(token, prefix, &g, logger)
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
