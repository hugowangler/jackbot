package bot

import (
	"fmt"
	"jackbot/internal/game"
	"strings"
)

func HandleInput(input string, g *game.Game) string {
	switch {
	case strings.HasPrefix(input, "!join"):
		row, err := handleJoin(input, g)
		if err != nil {
			return err.Error()
		}
		return row.Format()
	}
	return ""
}

func handleJoin(msg string, g *game.Game) (game.Row, error) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!join"))

	var row game.Row
	var err error

	switch msg {
	case "":
		return game.Row{}, fmt.Errorf("help me")
	case "random":
		row = g.GetRandomRow()
	default:
		row, err = g.ParseRow(msg)
		if err != nil {
			return game.Row{}, err
		}
	}

	return row, nil
}
