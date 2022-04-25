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

const okResponse = "game created: jacken"
const validCreateGame = "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 " +
	"--entryfee=5"

type CommandsTS struct {
	suite.Suite
	testDbContainer      testcontainers.Container
	db                   *gorm.DB
	mockedCommandHandler *CommandHandler
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
		name     string
		input    string
		exp      string
		authorId string
	}{
		{
			name:     "valid",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      okResponse,
			authorId: "178632146762596352",
		},
		{
			name:     "missing name",
			input:    "!creategame --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      "name is required",
			authorId: "178632146762596352",
		},
		{
			name:     "missing numbers",
			input:    "!creategame --name=jacken --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      "numbers is required",
			authorId: "178632146762596352",
		},
		{
			name:     "missing numbersrange",
			input:    "!creategame --name=jacken --numbers=5 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      "numbersrange is required",
			authorId: "178632146762596352",
		},
		{
			name:     "missing bonusnumbers",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusrange=12 --entryfee=5",
			exp:      "bonusnumbers is required",
			authorId: "178632146762596352",
		},
		{
			name:     "missing bonusrange",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --entryfee=5",
			exp:      "bonusrange is required",
			authorId: "178632146762596352",
		},
		{
			name:     "missing entryfee",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12",
			exp:      "entryfee is required",
			authorId: "178632146762596352",
		},
		{
			name:     "numbers not an integer",
			input:    "!creategame --name=jacken --numbers=a --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      "numbers must be an integer",
			authorId: "178632146762596352",
		},
		{
			name:     "numbersrange not an integer",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=b --bonusnumbers=2 --bonusrange=12 --entryfee=5",
			exp:      "numbersrange must be an integer",
			authorId: "178632146762596352",
		},
		{
			name:     "bonusnumbers not an integer",
			input:    "!creategame  --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=a --bonusrange=12 --entryfee=5",
			exp:      "bonusnumbers must be an integer",
			authorId: "178632146762596352",
		},
		{
			name:     "bonusrange not an integer",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=e --entryfee=5",
			exp:      "bonusrange must be an integer",
			authorId: "178632146762596352",
		},
		{
			name:     "entryfee not an integer",
			input:    "!creategame --name=jacken --numbers=5 --numbersrange=50 --bonusnumbers=2 --bonusrange=12 --entryfee=a",
			exp:      "entryfee must be an integer",
			authorId: "178632146762596352",
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
	err := test.SeedGame(&test.MockGame, s.db)
	assert.Nil(s.T(), err)

	res := s.mockedCommandHandler.HandleInput(validCreateGame, "178632146762596352")
	assert.Equal(s.T(), "game with name jacken already exists", res)
}
