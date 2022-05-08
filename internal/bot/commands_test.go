package bot

import (
	"context"
	"jackbot/db"
	"jackbot/internal/row"
	"jackbot/internal/utils"
	"jackbot/test"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

const okResponse = "game created: jacken"
const validCreateGame = "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 " +
	"--entryfee=5 --accountant=abc123"

type CommandsTS struct {
	suite.Suite
	testDbContainer      testcontainers.Container
	db                   *gorm.DB
	mockedCommandHandler *CommandHandler
	logger               *zap.SugaredLogger
}

func (s *CommandsTS) SetupSuite() {
	s.testDbContainer = test.StartTestDb(&s.Suite)
	host, err := s.testDbContainer.Host(context.Background())
	if err != nil {
		panic(err)
	}
	port, err := s.testDbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		panic(err)
	}

	os.Setenv("POSTGRES_HOST", host)
	os.Setenv("POSTGRES_PORT", port.Port())
	os.Setenv("POSTGRES_DB", "test")
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")

	dbConn, err := db.NewConn()
	if err != nil {
		panic(err)
	}
	s.db = dbConn

	err = test.RunMigrations(s.db)
	if err != nil {
		panic(err)
	}

	s.logger = utils.NewLogger()

	s.mockedCommandHandler = &CommandHandler{
		db:     s.db,
		logger: s.logger,
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
	var err error

	err = test.SeedUser(&test.MockUser, s.db)
	assert.Nil(s.T(), err)

	err = test.SeedPermission(&test.MockPermission, s.db)
	assert.Nil(s.T(), err)

	tests := []struct {
		name     string
		input    string
		exp      string
		authorId string
	}{
		{
			name:     "valid",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      okResponse,
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing name",
			input:    "!creategame --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "name is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing numbers",
			input:    "!creategame --name=jacken --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "numbers is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing numbersrange",
			input:    "!creategame --name=jacken --numbers=5 --bonusnumbers=2 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "numbersrange is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing bonusnumbers",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "bonusnumbers is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing bonusrange",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --entryfee=5 --accountant=abc123",
			exp:      "bonusrange is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing entryfee",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --accountant=abc123",
			exp:      "entryfee is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "missing accountant",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      "accountant is required",
			authorId: test.MockUser.Id,
		},
		{
			name:     "numbers not an integer",
			input:    "!creategame --name=jacken --numbers=a --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "numbers must be an integer",
			authorId: test.MockUser.Id,
		},
		{
			name:     "numbersrange not an integer",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=b --bonusnumbers=2 --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "numbersrange must be an integer",
			authorId: test.MockUser.Id,
		},
		{
			name:     "bonusnumbers not an integer",
			input:    "!creategame  --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=a --bonusrange=12 --entryfee=5 --accountant=abc123",
			exp:      "bonusnumbers must be an integer",
			authorId: test.MockUser.Id,
		},
		{
			name:     "bonusrange not an integer",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=e --entryfee=5 --accountant=abc123",
			exp:      "bonusrange must be an integer",
			authorId: test.MockUser.Id,
		},
		{
			name:     "entryfee not an integer",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=a --accountant=abc123",
			exp:      "entryfee must be an integer",
			authorId: test.MockUser.Id,
		},
	}
	for _, tt := range tests {
		s.T().Run(
			tt.name, func(t *testing.T) {
				actual := s.mockedCommandHandler.HandleInput(tt.input, tt.authorId)
				assert.Equal(t, tt.exp, actual)
			},
		)
	}
}

func (s *CommandsTS) TestCommands_HandleInput_CreateGame_NameAlreadyExists() {
	var err error

	err = test.SeedUser(&test.MockUser, s.db)
	assert.Nil(s.T(), err)

	err = test.SeedGame(&test.MockGame, s.db)
	assert.Nil(s.T(), err)

	err = test.SeedPermission(&test.MockPermission, s.db)
	assert.Nil(s.T(), err)

	res := s.mockedCommandHandler.HandleInput(validCreateGame, test.MockUser.Id)
	assert.Equal(s.T(), "game with name jacken already exists", res)
}

func (s *CommandsTS) TestCommands_NewCmdHandler() {
	cmd := NewCmdHandler(s.db, s.logger)
	assert.NotNil(s.T(), cmd)
	assert.Nil(s.T(), cmd.rowHandler)
	assert.Nil(s.T(), cmd.game)
}

func (s *CommandsTS) TestCommands_NewCmdHandler_WithGame() {
	cmd := NewCmdHandler(s.db, s.logger, WithGame(&test.MockGame))
	assert.NotNil(s.T(), cmd)
	assert.Nil(s.T(), cmd.rowHandler)
	assert.NotNil(s.T(), cmd.game)
}

func (s *CommandsTS) TestCommands_NewCmdHandler_WithRowHandler() {
	cmd := NewCmdHandler(s.db, s.logger, WithRowHandler(&row.Handler{}))
	assert.NotNil(s.T(), cmd)
	assert.NotNil(s.T(), cmd.rowHandler)
	assert.Nil(s.T(), cmd.game)
}
