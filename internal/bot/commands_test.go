package bot

import (
	"context"
	"jackbot/db"
	"jackbot/internal/utils"
	"jackbot/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

type CommandsTS struct {
	suite.Suite
	testDbContainer      testcontainers.Container
	db                   *gorm.DB
	mockedCommandHandler *CommandHandler
}

func (s *CommandsTS) SetupSuite() {
	s.testDbContainer = test.StartTestDb(&s.Suite)
	host, err := s.testDbContainer.ContainerIP(context.Background())
	if err != nil {
		panic(err)
	}

	os.Setenv("POSTGRES_HOST", host)
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "test")
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")

	db, err := db.NewConn()
	if err != nil {
		panic(err)
	}
	s.db = db

	err = test.RunMigrations(s.db)
	if err != nil {
		panic(err)
	}

	logger := utils.NewLogger()

	s.mockedCommandHandler = &CommandHandler{
		db:     s.db,
		logger: logger,
	}
}

func (s *CommandsTS) TearDownTest() {
	err := test.CleanTestDb(s.db)
	if err != nil {
		panic(err)
	}
}

func TestCommandsTS(t *testing.T) {
	suite.Run(t, new(CommandsTS))
}

func (s *CommandsTS) TestCommands_HandleInput_CreateGame() {
	tests := []struct {
		name  string
		input string
		exp   string
	}{
		{
			name:  "valid create game command",
			input: "!createGame --name=Jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:   "game created: jacken",
		},
	}
	for _, tt := range tests {
		s.T().Run(
			tt.name, func(t *testing.T) {
				actual := s.mockedCommandHandler.HandleInput(tt.input)
				assert.Equal(t, tt.exp, actual)
			},
		)
	}
}
