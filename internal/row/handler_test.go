package row

import (
	"jackbot/db/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockRowHandler = Handler{
	&models.Game{
		Numbers:      5,
		NumbersRange: 10,
		BonusNumbers: 2,
		BonusRange:   5,
	},
}

func TestHandler_ParseRow(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		exp    models.Row
		expErr bool
	}{
		{
			name:  "valid row",
			input: "10 2 3 4 5 - 1 5",
			exp: models.Row{
				Numbers:      []int32{10, 2, 3, 4, 5},
				BonusNumbers: []int32{1, 5},
			},
			expErr: false,
		},
		{
			name:  "no space between separator",
			input: "1 2 3 4 5-1 2",
			exp: models.Row{
				Numbers:      []int32{1, 2, 3, 4, 5},
				BonusNumbers: []int32{1, 2},
			},
			expErr: false,
		},
		{
			name:   "wrong separator",
			input:  "1 2 3 4 5 , 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "Numbers not a number",
			input:  "1 2 aabbcc 4 5 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "BonusNumbers Numbers not a number",
			input:  "1 2 3 4 5 - 1 asdad",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "missing separator",
			input:  "1 2 3 4 5 1",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "too few Numbers",
			input:  "1 2 3 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "too few BonusNumbers Numbers",
			input:  "1 2 3 4 5 - 1",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "too many Numbers",
			input:  "1 2 3 4 5 6 7 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "too many BonusNumbers Numbers",
			input:  "1 2 3 4 5 - 1 2 3 4",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "Numbers repeated number",
			input:  "1 2 3 1 5 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "BonusNumbers Numbers repeated number",
			input:  "1 2 3 4 5 - 2 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "Numbers out of range",
			input:  "1 2 33 4 5 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "BonusNumbers Numbers out of range",
			input:  "1 2 3 4 5 - 103 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "number is zero",
			input:  "1 2 0 4 5 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "BonusNumbers Numbers is zero",
			input:  "1 2 3 4 5 - 0 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "number is negative",
			input:  "1 2 3 4 -5 - 1 2",
			exp:    models.Row{},
			expErr: true,
		},
		{
			name:   "BonusNumbers Numbers is negative",
			input:  "1 2 3 4 5 - 1 -2",
			exp:    models.Row{},
			expErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual, err := mockRowHandler.ParseRow(tt.input)
				if (err != nil) != tt.expErr {
					t.Errorf("Row.ParseRow() error = %v, expError=%v", err, tt.expErr)
				}
				assert.Equal(t, tt.exp, actual)
			},
		)
	}
}
