package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const ballSeparator = "-"

type Game struct {
	Balls      int
	BallsRange int
	Bonus      int
	BonusRange int
}

type Row struct {
	Balls []int
	Bonus []int
}

func (r *Row) Format() string {
	var res string
	for _, b := range r.Balls {
		res += fmt.Sprintf("%d ", b)
	}
	res += fmt.Sprintf("%s ", ballSeparator)
	for _, b := range r.Bonus {
		res += fmt.Sprintf("%d ", b)
	}
	return res
}

type WrongInputFormatError struct {
}

func (e *WrongInputFormatError) Error() string {
	return fmt.Sprintf(
		"wrong input format for the !join command:\n correct format is !join <balls> %s <bonus balls>. Where "+
			"<balls> and <bonus balls> are space seperated numbers in the balls and bonus balls number range",
		ballSeparator,
	)
}

func (g *Game) GetRandomRow() Row {
	var row Row
	selectedBalls := make(map[int]struct{}, g.Balls)
	for i := 0; i < g.Balls; i++ {
		for {
			ball := rand.Intn(g.BallsRange) + 1
			if _, exists := selectedBalls[ball]; !exists {
				selectedBalls[ball] = struct{}{}
				row.Balls = append(row.Balls, ball)
				break
			}
		}
	}
	selectedBonus := make(map[int]struct{}, g.Bonus)
	for i := 0; i < g.Bonus; i++ {
		for {
			ball := rand.Intn(g.BonusRange) + 1
			if _, exists := selectedBonus[ball]; !exists {
				selectedBonus[ball] = struct{}{}
				row.Bonus = append(row.Bonus, ball)
				break
			}
		}
	}
	return row
}

func (g *Game) ParseRow(input string) (Row, error) {
	var row Row
	var err error
	// "!join 1 2 3 4 5 - 7 8"
	split := strings.Split(input, ballSeparator)
	if len(split) != 2 {
		return Row{}, &WrongInputFormatError{}
	}
	row.Balls, err = parseBalls(strings.TrimSpace(split[0]), g.Balls, g.BallsRange)
	if err != nil {
		return Row{}, err
	}
	row.Bonus, err = parseBalls(strings.TrimSpace(split[1]), g.Bonus, g.BonusRange)
	if err != nil {
		return Row{}, err
	}
	return row, nil
}

func parseBalls(input string, numBalls int, ballRange int) ([]int, error) {
	var balls []int
	strBalls := strings.Split(input, " ")
	if len(strBalls) != numBalls {
		return []int{}, &WrongInputFormatError{}
	}
	selectedBalls := make(map[int]struct{}, numBalls)
	for _, b := range strBalls {
		i, err := strconv.Atoi(b)
		if err != nil {
			return []int{}, &WrongInputFormatError{}
		}
		if i <= 0 || i > ballRange {
			return []int{}, &WrongInputFormatError{}
		}
		if _, exists := selectedBalls[i]; exists {
			return []int{}, &WrongInputFormatError{}
		}
		selectedBalls[i] = struct{}{}
		balls = append(balls, i)
	}
	return balls, nil
}
