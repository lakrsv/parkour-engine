package main

import (
	"reflect"
	"unicode"

	"github.com/lakrsv/parkour-engine/engine"
)

func ButtonBlueprint(x int, y int, char rune) []any {
	return []any{PositionComponent{X: x, Y: y},
		RenderComponent{Character: Button},
		TriggerComponent{Symbol: unicode.ToUpper(char)}}
}

func ClosedDoorBlueprint(x int, y int, char rune) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		DeferDoorRenderComponent{},
		ObstacleComponent{},
		TriggeredComponent{Symbol: char, Action: func(entity uint32, w *engine.World) {
			doorOpenPlayCount := reflect.ValueOf(w.GetUniqueComponent(reflect.TypeOf(DoorOpenPlayCountComponent{}))).Interface().(DoorOpenPlayCountComponent)
			playDoorOpenSound(doorOpenPlayCount.Count)
			w.ReplaceUniqueComponent(DoorOpenPlayCountComponent{Count: doorOpenPlayCount.Count + 1})

			w.RemoveComponent(entity, reflect.TypeOf(ObstacleComponent{}))
			w.ReplaceComponent(entity, RenderComponent{Character: OpenDoor})
			w.AddComponent(entity, FloorComponent{})
		}}}
}

func OpenDoorBlueprint(x int, y int, char rune) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: OpenDoor},
		FloorComponent{},
		TriggeredComponent{Symbol: char, Action: func(entity uint32, w *engine.World) {
			doorOpenPlayCount := reflect.ValueOf(w.GetUniqueComponent(reflect.TypeOf(DoorOpenPlayCountComponent{}))).Interface().(DoorOpenPlayCountComponent)
			playDoorOpenSound(doorOpenPlayCount.Count)
			w.ReplaceUniqueComponent(DoorOpenPlayCountComponent{Count: doorOpenPlayCount.Count + 1})

			w.RemoveComponent(entity, reflect.TypeOf(RenderComponent{}))
			w.RemoveComponent(entity, reflect.TypeOf(FloorComponent{}))
			w.AddComponent(entity, DeferDoorRenderComponent{})
			w.AddComponent(entity, ObstacleComponent{})
		}}}
}

func WallBlueprint(x int, y int) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: Wall},
		ObstacleComponent{},
	}
}

func FloorBlueprint(x int, y int) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: Floor},
		FloorComponent{},
	}
}

func PlayerBlueprint(x int, y int) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: Player},
		PlayerInputComponent{},
		MoveComponent{},
		FacingComponent{},
		InteractsWithTriggersComponent{color: struct{ R, G, B uint8 }{R: 0, G: 255, B: 0}},
		SummonComponent{color: struct{ R, G, B uint8 }{R: 0, G: 255, B: 255}},
	}
}

func ExitBlueprint(x int, y int, level int) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: Exit},
		TriggerComponent{Symbol: Exit},
		TriggeredComponent{Symbol: Exit, Action: func(entity uint32, w *engine.World) {
			playGoalSound()
			if err := w.Close(); err != nil {
				panic(err)
			}
			Run(level)
		}},
	}
}

func SummonBlueprint(x int, y int, facingX int, facingY int, color Color) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: Summon},
		SummonInputComponent{X: facingX, Y: facingY},
		MoveComponent{},
		ColorComponent{color: color},
		InteractsWithTriggersComponent{color: color},
	}
}

func SummonPickupBlueprint(x int, y int, color Color) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: SummonPickup},
		ColorComponent{color: color},
		SummonPickupComponent{color: color},
	}
}
