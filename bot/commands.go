package bot

import (
	"fmt"
	"jackbot/game"
	"strings"
)

func handleJoin(msg string, g *game.Game) (game.Row, string) {
	msg = strings.TrimSpace(strings.TrimPrefix(msg, "!join"))

	var row game.Row
	var err error

	switch msg {
	case "":
		return game.Row{}, "help me"
	case "random":
		row = g.GetRandomRow()
	default:
		row, err = g.ParseRow(msg)
		if err != nil {
			return game.Row{}, fmt.Sprintf("ERROR: %s", err.Error())
		}
	}

	return row, ""
}
