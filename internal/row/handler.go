package row

import (
	"fmt"
	"jackbot/db/models"
	"math/rand"
	"strconv"
	"strings"
)

const separator = "-"

type Handler struct {
	models.Game
}

type InvalidRowError struct {
}

func (e *InvalidRowError) Error() string {
	return fmt.Sprintf(
		"wrong input format for the !join command:\n correct format is !join <numbers> %s <bonus numbers>. Where "+
			"<numbers> and <bonus numbers> are space seperated numbers in the numbers and bonus numbers number range",
		separator,
	)
}

func (r *Handler) GetRandomRow() models.Row {
	var row models.Row
	selectedNumbers := make(map[int]struct{}, r.Numbers)
	for i := 0; i < r.Numbers; i++ {
		for {
			number := rand.Intn(r.NumbersRange) + 1
			if _, exists := selectedNumbers[number]; !exists {
				selectedNumbers[number] = struct{}{}
				row.Numbers = append(row.Numbers, int32(number))
				break
			}
		}
	}
	selectedBonusNumbers := make(map[int]struct{}, r.BonusNumbers)
	for i := 0; i < r.BonusNumbers; i++ {
		for {
			number := rand.Intn(r.BonusRange) + 1
			if _, exists := selectedBonusNumbers[number]; !exists {
				selectedBonusNumbers[number] = struct{}{}
				row.BonusNumbers = append(row.Numbers, int32(number))
				break
			}
		}
	}
	return row
}

func (r *Handler) ParseRow(input string) (models.Row, error) {
	var row models.Row
	var err error
	// "!join 1 2 3 4 5 - 7 8"
	split := strings.Split(input, separator)
	if len(split) != 2 {
		return models.Row{}, &InvalidRowError{}
	}
	row.Numbers, err = parseNumbers(strings.TrimSpace(split[0]), r.Numbers, r.NumbersRange)
	if err != nil {
		return models.Row{}, err
	}
	row.BonusNumbers, err = parseNumbers(strings.TrimSpace(split[1]), r.BonusNumbers, r.BonusRange)
	if err != nil {
		return models.Row{}, err
	}
	return row, nil
}

func parseNumbers(input string, numNumbers int, numberRange int) ([]int32, error) {
	var numbers []int32
	strNumbers := strings.Split(input, " ")
	if len(strNumbers) != numNumbers {
		return []int32{}, &InvalidRowError{}
	}
	selectedNumbers := make(map[int]struct{}, numNumbers)
	for _, b := range strNumbers {
		i, err := strconv.Atoi(b)
		if err != nil {
			return []int32{}, &InvalidRowError{}
		}
		if i <= 0 || i > numberRange {
			return []int32{}, &InvalidRowError{}
		}
		if _, exists := selectedNumbers[i]; exists {
			return []int32{}, &InvalidRowError{}
		}
		selectedNumbers[i] = struct{}{}
		numbers = append(numbers, int32(i))
	}
	return numbers, nil
}
