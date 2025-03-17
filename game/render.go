package main

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

type Color struct {
	R, G, B uint8
}

type RunePalette struct {
	runeColors map[rune]Color
}

func NewRunePalette(runeColors map[rune]Color) RunePalette {
	return RunePalette{runeColors}
}

func (p RunePalette) GetColor(r rune) Color {
	if col, ok := p.runeColors[r]; ok {
		return col
	}
	return Color{R: 255, G: 255, B: 255} // White as default
}
