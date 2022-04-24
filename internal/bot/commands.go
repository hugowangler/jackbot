package bot

import (
	"errors"
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
}

func (c *CommandHandler) HandleInput(input string) string {
	switch {
	case strings.HasPrefix(input, "!join"):
		row, err := c.handleJoin(input)
		if err != nil {
			return err.Error()
		}
		return row.NumbersToString()
	case strings.HasPrefix(input, "!creategame"):
		game, err := c.handleCreateGame(input)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("game created: %s", game.Name)
	}
	return ""
}

func (c *CommandHandler) handleJoin(msg string) (models.Row, error) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!join"))

	var row models.Row
	var err error

	switch msg {
	case "":
		return models.Row{}, fmt.Errorf("help me")
	case "random":
		row = c.rowHandler.GetRandomRow()
	default:
		row, err = c.rowHandler.ParseRow(msg)
		if err != nil {
			return models.Row{}, err
		}
	}

	res := c.db.Create(&row)
	if res.Error != nil {
		err = utils.LogServerError(err, c.logger)
		return models.Row{}, err
	}

	return row, nil
}

func (c *CommandHandler) handleCreateGame(msg string) (models.Game, error) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!creategame"))

	var game models.Game
	var exists bool
	var err error

	args := ParseArguments(msg)

	c.logger.Debug("args", args)

	game.Name, exists = args["name"]
	if !exists {
		return models.Game{}, fmt.Errorf("name is required")
	}

	strNumbers, exists := args["numbers"]
	if !exists {
		return models.Game{}, fmt.Errorf("numbers is required")
	}
	c.logger.Debug(strNumbers)
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
		return models.Game{}, fmt.Errorf("strBonusRange is required")
	}
	game.BonusRange, err = strconv.Atoi(strBonusRange)
	if err != nil {
		return models.Game{}, fmt.Errorf("bonusRange must be an integer")
	}

	strEntryFee, exists := args["entryfee"]
	if !exists {
		return models.Game{}, fmt.Errorf("entryfee is required")
	}
	game.EntryFee, err = strconv.Atoi(strEntryFee)
	if err != nil {
		return models.Game{}, fmt.Errorf("entryfee must be an integer")
	}

	game.Active = true

	err = models.CreateGame(&game, c.db)
	if err != nil {
		if errors.Is(err, &models.GameAlreadyExistsError{}) {
			return models.Game{}, err
		}

		err = utils.LogServerError(err, c.logger)
		return models.Game{}, err
	}

	return game, nil
}
