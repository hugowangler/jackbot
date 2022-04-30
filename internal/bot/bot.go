package bot

import (
	"jackbot/db/models"
	"jackbot/internal/row"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bot struct {
	logger     *zap.SugaredLogger
	game       *models.Game
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

	response := b.cmdHandler.HandleInput(msg, m.Author.ID)
	if response != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			b.logger.With("error", err).Warn("failed to respond")
		}
	}
}

func NewBot(token string, prefix string, game *models.Game, logger *zap.SugaredLogger, db *gorm.DB) *Bot {
	bot := &Bot{
		token:      token,
		prefix:     prefix,
		logger:     logger,
		game:       nil,
		cmdHandler: nil,
	}
	if game != nil {
		bot.game = game
		bot.cmdHandler = NewCmdHandler(db, logger, WithGame(game), WithRowHandler(&row.Handler{Game: game}))
	} else {
		bot.cmdHandler = NewCmdHandler(db, logger)
	}
	return bot
}
