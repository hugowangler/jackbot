package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockGame = Game{
	Balls:      5,
	BallsRange: 10,
	Bonus:      2,
	BonusRange: 5,
}

func TestGame_ParseRow(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		exp    Row
		expErr bool
	}{
		{
			name:  "valid row",
			input: "10 2 3 4 5 - 1 5",
			exp: Row{
				Balls: []int{10, 2, 3, 4, 5},
				Bonus: []int{1, 5},
			},
			expErr: false,
		},
		{
			name:  "no space between separator",
			input: "1 2 3 4 5-1 2",
			exp: Row{
				Balls: []int{1, 2, 3, 4, 5},
				Bonus: []int{1, 2},
			},
			expErr: false,
		},
		{
			name:   "wrong separator",
			input:  "1 2 3 4 5 , 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "balls not a number",
			input:  "1 2 aabbcc 4 5 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "bonus balls not a number",
			input:  "1 2 3 4 5 - 1 asdad",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "missing separator",
			input:  "1 2 3 4 5 1",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "too few balls",
			input:  "1 2 3 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "too few bonus balls",
			input:  "1 2 3 4 5 - 1",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "too many balls",
			input:  "1 2 3 4 5 6 7 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "too many bonus balls",
			input:  "1 2 3 4 5 - 1 2 3 4",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "balls repeated number",
			input:  "1 2 3 1 5 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "bonus balls repeated number",
			input:  "1 2 3 4 5 - 2 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "balls out of range",
			input:  "1 2 33 4 5 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "bonus balls out of range",
			input:  "1 2 3 4 5 - 103 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "ball is zero",
			input:  "1 2 0 4 5 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "bonus balls is zero",
			input:  "1 2 3 4 5 - 0 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "ball is negative",
			input:  "1 2 3 4 -5 - 1 2",
			exp:    Row{},
			expErr: true,
		},
		{
			name:   "bonus balls is negative",
			input:  "1 2 3 4 5 - 1 -2",
			exp:    Row{},
			expErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual, err := mockGame.ParseRow(tt.input)
				if (err != nil) != tt.expErr {
					t.Errorf("Game.ParseRow() error = %v, expError=%v", err, tt.expErr)
				}
				assert.Equal(t, tt.exp, actual)
			},
		)
	}
}
