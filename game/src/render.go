package src

import (
	"github.com/fatih/color"
)

const (
	Floor           = ' '
	OpenDoor        = '\''
	Wall            = '#'
	Button          = '='
	TriggeredButton = '_'
	Player          = '@'
	Summon          = 'S'
	DoorHorizontal  = '-'
	DoorVertical    = '|'
	LeftIndicator   = '˂'
	RightIndicator  = '˃'
	UpIndicator     = '˄'
	DownIndicator   = '˅'
	Exit            = '%'
	SummonPickup    = '~'
	CyanSummon      = '0'
	RedSummon       = '1'
	YellowSummon    = '2'
)

type RunePalette struct {
	runeColors map[rune]color.Color
}

func NewRunePalette(runeColors map[rune]color.Color) RunePalette {
	return RunePalette{runeColors}
}

func (p RunePalette) GetColor(r rune) *color.Color {
	if col, ok := p.runeColors[r]; ok {
		return &col
	}
	return color.New(color.FgWhite)
}
