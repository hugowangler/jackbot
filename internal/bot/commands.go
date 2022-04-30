package bot

import (
	"fmt"
	"jackbot/db/models"
	"jackbot/internal/row"
	"jackbot/internal/utils"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommandHandler struct {
	db         *gorm.DB
	logger     *zap.SugaredLogger
	rowHandler *row.Handler
	game       *models.Game
}

func (c *CommandHandler) HandleInput(input string, authorId string) string {
	switch {
	case strings.HasPrefix(input, "!join"):
		dbRow, err := c.handleJoin(input, authorId)
		if err != nil {
			return err.Error()
		}

		accountant, err := models.GetAccountant(c.game, c.db)
		if err != nil {
			return err.Error()
		}

		returnMessage := fmt.Sprintf("Your row has been registered: `%s - %s`. Please send the entry fee of %d SEK to %s.",
			dbRow.NumbersToString(),
			dbRow.BonusNumbersToString(),
			c.game.EntryFee,
			accountant.Mobile)

		if dbRow.User.Mobile == "" || dbRow.User.Name == "" {
			returnMessage = fmt.Sprintf(`%s

To be able to win, we need some additional information to be able to confim your transactions.

To do this, please type !setup <name> <mobile>
example: !setup --name=John Doe --mobile=0731234567

Note: It's important that the mobile number and name that you enter matches your Swish information exactly. If you are unsure, please contact a server administrator.`,
				returnMessage,
			)
		}

		return returnMessage
	case strings.HasPrefix(input, "!createraffle"):
		raffle, err := c.handleCreateRaffle(input)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("Raffle created: %d", raffle.GameId)
	case strings.HasPrefix(input, "!creategame"):
		hasPermission, err := models.HasPermissions(authorId, []int{models.MasterAdmin}, c.db)
		if err != nil {
			return utils.LogServerError(err, c.logger).Error()
		}
		if !hasPermission {
			return "you don't have permission to do this action"
		}

		game, err := c.handleCreateGame(input)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("game created: %s", game.Name)
	case strings.HasPrefix(input, "!setup"):
		user, err := c.handleSetupUser(input, authorId)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("user created: %s", user.Name)
	}
	return ""
}

func (c *CommandHandler) handleJoin(msg string, userId string) (models.Row, error) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!join"))

	var dbRow models.Row
	var err error

	switch msg {
	case "":
		return models.Row{}, fmt.Errorf("help me")
	case "random":
		dbRow = c.rowHandler.GetRandomRow()
	default:
		dbRow, err = c.rowHandler.ParseRow(msg)
		if err != nil {
			return models.Row{}, err
		}
	}

	user, err := models.GetUser(userId, c.db)
	if err != nil {
		err = utils.LogServerError(err, c.logger)
		return models.Row{}, err
	}
	if user == nil {
		user = &models.User{
			Id: userId,
		}
		if models.CreateUser(*user, c.db) != nil {
			if err, ok := err.(*models.UserAlreadyExistsError); ok {
				return models.Row{}, err
			}

			err = utils.LogServerError(err, c.logger)
			return models.Row{}, err
		}
	}
	dbRow.UserId = user.Id
	dbRow.User = *user

	raffle, err := models.GetActiveRaffle(c.game.Id, c.db)
	if err != nil {
		if err, ok := err.(*models.PreviousRaffleNotCompletedError); ok {
			return models.Row{}, err
		}

		err = utils.LogServerError(err, c.logger)
		return models.Row{}, err
	}
	dbRow.RaffleId = raffle.Id

	res := c.db.Create(&dbRow)
	if res.Error != nil {
		err = utils.LogServerError(err, c.logger)
		return models.Row{}, err
	}

	return dbRow, nil
}

func (c *CommandHandler) handleCreateRaffle(msg string) (models.Raffle, error) {
	var raffle models.Raffle

	raffle.GameId = c.game.Id

	err := models.CreateRaffle(&raffle, c.db)

	if err != nil {
		if err, ok := err.(*models.PreviousRaffleNotCompletedError); ok {
			return models.Raffle{}, err
		}

		err = utils.LogServerError(err, c.logger)
		return models.Raffle{}, err
	}

	return raffle, nil
}

func (c *CommandHandler) handleCreateGame(msg string) (models.Game, error) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!creategame"))

	var game models.Game
	var exists bool
	var err error

	args := ParseArguments(msg)

	game.Name, exists = args["name"]
	if !exists {
		return models.Game{}, fmt.Errorf("name is required")
	}

	strNumbers, exists := args["numbers"]
	if !exists {
		return models.Game{}, fmt.Errorf("numbers is required")
	}
	game.Numbers, err = strconv.Atoi(strNumbers)
	if err != nil {
		return models.Game{}, fmt.Errorf("numbers must be an integer")
	}

	strNumbersRange, exists := args["numbersrange"]
	if !exists {
		return models.Game{}, fmt.Errorf("numbersrange is required")
	}
	game.NumbersRange, err = strconv.Atoi(strNumbersRange)
	if err != nil {
		return models.Game{}, fmt.Errorf("numbersrange must be an integer")
	}

	strBonusNumbers, exists := args["bonusnumbers"]
	if !exists {
		return models.Game{}, fmt.Errorf("bonusnumbers is required")
	}
	game.BonusNumbers, err = strconv.Atoi(strBonusNumbers)
	if err != nil {
		return models.Game{}, fmt.Errorf("bonusnumbers must be an integer")
	}

	strBonusRange, exists := args["bonusrange"]
	if !exists {
		return models.Game{}, fmt.Errorf("bonusrange is required")
	}
	game.BonusRange, err = strconv.Atoi(strBonusRange)
	if err != nil {
		return models.Game{}, fmt.Errorf("bonusrange must be an integer")
	}

	strEntryFee, exists := args["entryfee"]
	if !exists {
		return models.Game{}, fmt.Errorf("entryfee is required")
	}
	game.EntryFee, err = strconv.Atoi(strEntryFee)
	if err != nil {
		return models.Game{}, fmt.Errorf("entryfee must be an integer")
	}

	game.AccountantId, exists = args["accountant"]
	if !exists {
		return models.Game{}, fmt.Errorf("accountant is required")
	}

	game.Active = true

	err = models.CreateGame(&game, c.db)
	if err != nil {
		if err, ok := err.(*models.GameAlreadyExistsError); ok {
			return models.Game{}, err
		}

		err = utils.LogServerError(err, c.logger)
		return models.Game{}, err
	}

	return game, nil
}

func (c *CommandHandler) handleSetupUser(msg string, userId string) (models.User, error) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!setup"))

	var err error
	var exists bool

	user, err := models.GetUser(userId, c.db)
	if err != nil {
		err = utils.LogServerError(err, c.logger)
		return models.User{}, err
	}
	if user == nil {
		user = &models.User{
			Id: userId,
		}
		if models.CreateUser(*user, c.db) != nil {
			if err, ok := err.(*models.UserAlreadyExistsError); ok {
				return models.User{}, err
			}

			err = utils.LogServerError(err, c.logger)
			return models.User{}, err
		}
	}

	args := ParseArguments(msg)

	user.Name, exists = args["name"]
	if !exists {
		return models.User{}, fmt.Errorf("name is required")
	}

	user.Mobile, exists = args["mobile"]
	if !exists {
		return models.User{}, fmt.Errorf("mobile is required")
	}

	if c.db.Save(&user).Error != nil {
		err = utils.LogServerError(err, c.logger)
		return models.User{}, err
	}

	return *user, nil
}

func NewCmdHandler(
	db *gorm.DB,
	logger *zap.SugaredLogger,
	opts ...func(c *CommandHandler),
) *CommandHandler {
	cmd := &CommandHandler{
		db:     db,
		logger: logger,
	}
	for _, o := range opts {
		o(cmd)
	}
	return cmd
}

func WithGame(game *models.Game) func(c *CommandHandler) {
	return func(c *CommandHandler) {
		c.game = game
	}
}

func WithRowHandler(handler *row.Handler) func(c *CommandHandler) {
	return func(c *CommandHandler) {
		c.rowHandler = handler
	}
}
