package bot

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"jackbot/internal/game"
	"strings"
)

type Bot struct {
	logger  *zap.SugaredLogger
	game    *game.Game
	userId  string
	session *discordgo.Session
	token   string
	prefix  string
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

	var response string
	switch {
	case strings.HasPrefix(msg, "!join"):
		row, errMsg := handleJoin(msg, b.game)
		if errMsg != "" {
			response += errMsg
		}
		response += row.Format()
	}

	if response != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			b.logger.With("error", err).Warn("failed to respond")
		}
	}
}

func NewBot(token string, prefix string, game *game.Game, logger *zap.SugaredLogger) *Bot {
	return &Bot{
		token:  token,
		prefix: prefix,
		game:   game,
		logger: logger,
	}
}
