package bot

import (
	"jackbot/db/models"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bot struct {
	logger     *zap.SugaredLogger
	game       models.Game
	userId     string
	session    *discordgo.Session
	token      string
	prefix     string
	cmdHandler *CommandHandler
}

func (b *Bot) Start() error {
	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return err
	}

	u, err := session.User("@me")
	if err != nil {
		return err
	}
	b.userId = u.ID
	b.session = session

	session.AddHandler(b.messageHandler)

	err = session.Open()
	if err != nil {
		return err
	}
	b.logger.Info("jackbot is running")
	return nil
}

func (b *Bot) Close() error {
	err := b.session.Close()
	return err
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.userId {
		return
	}

	msg := strings.TrimSpace(strings.ToLower(m.Content))

	if !strings.HasPrefix(msg, b.prefix) {
		return
	}

	response := b.cmdHandler.HandleInput(msg)
	if response != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			b.logger.With("error", err).Warn("failed to respond")
		}
	}
}

func NewBot(token string, prefix string, games []models.Game, logger *zap.SugaredLogger, db *gorm.DB) *Bot {
	bot := &Bot{
		token:  token,
		prefix: prefix,
		logger: logger,
		cmdHandler: &CommandHandler{
			db:     db,
			logger: logger,
		},
	}
	if len(games) > 0 {
		bot.game = games[0]
		bot.cmdHandler.rowHandler.Game = games[0]
	}

	return bot
}
