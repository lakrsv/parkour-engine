package main

import (
	"github.com/lakrsv/parkour-engine/engine"
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

type RenderComponent struct {
	Character rune
}

type SummonComponent struct {
	colorAttr int
	color     struct {
		R, G, B uint8
	}
}

type CreateSummonComponent struct{}

type ColorComponent struct {
	colorAttr int
	color     struct {
		R, G, B uint8
	}
}

type SummonPickupComponent struct {
	colorAttr int
	color     struct {
		R, G, B uint8
	}
}

type FloorComponent struct{}

type PositionComponent struct {
	X, Y int
}

type MoveComponent struct {
	X, Y int
}

type InteractsWithTriggersComponent struct {
	colorAttr int
	color     struct {
		R, G, B uint8
	}
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

func (g GridComponent) GetCell(x, y int) int {
	return y*g.Width + x
}

func (g GridComponent) GetPosition(cell int) (x int, y int) {
	x = cell % g.Width
	y = cell / g.Width
	return x, y
}
