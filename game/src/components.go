package src

import (
	"github.com/fatih/color"
	"github.com/lakrsv/parkour/engine"
)

type PlayerInputComponent struct {
}

type SummonInputComponent struct {
	X, Y int
}

type LevelComponent struct {
	Level  int
	Header []string
}

type DoorOpenPlayCountComponent struct {
	Count int
}

type InputComponent struct {
	keyState map[rune]bool
}

type RenderComponent struct {
	Character rune
}

type SummonComponent struct {
	colorAttr color.Attribute
	color     color.Color
}

type CreateSummonComponent struct{}

type ColorComponent struct {
	colorAttr color.Attribute
	color     color.Color
}

type SummonPickupComponent struct {
	colorAttr color.Attribute
	color     color.Color
}

type FloorComponent struct{}

type PositionComponent struct {
	X, Y int
}

type MoveComponent struct {
	X, Y int
}

type InteractsWithTriggersComponent struct {
	colorAttr color.Attribute
	color     color.Color
}

type FacingComponent struct {
	X, Y int
}

type DeferDoorRenderComponent struct {
}

type TriggerComponent struct {
	Symbol    rune
	Triggered bool
}

type TriggeredComponent struct {
	Symbol rune
	Action func(entity uint32, w *engine.World)
}

type GridComponent struct {
	Width, Height      int
	BackgroundEntities []uint32
	ForegroundEntities []uint32
	EffectEntities     []uint32
}

type ObstacleComponent struct {
}

func (c InputComponent) HasKey(key rune) bool {
	if state, ok := c.keyState[key]; ok {
		return state
	}
	return false
}

func (g GridComponent) GetCell(x, y int) int {
	return y*g.Width + x
}

func (g GridComponent) GetPosition(cell int) (x int, y int) {
	x = cell % g.Width
	y = cell / g.Width
	return x, y
}
