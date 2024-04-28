package game

import (
	"github.com/fatih/color"
	"github.com/lakrsv/parkour-engine/engine"
	"reflect"
	"unicode"
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
		InteractsWithTriggersComponent{color: *color.New(color.FgGreen), colorAttr: color.FgGreen},
		SummonComponent{color: *color.New(color.FgCyan), colorAttr: color.FgCyan},
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

func SummonBlueprint(x int, y int, facingX int, facingY int, colorAttr color.Attribute) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: Summon},
		SummonInputComponent{X: facingX, Y: facingY},
		MoveComponent{},
		ColorComponent{color: *color.New(colorAttr), colorAttr: colorAttr},
		InteractsWithTriggersComponent{color: *color.New(colorAttr), colorAttr: colorAttr},
	}
}

func SummonPickupBlueprint(x int, y int, colorAttr color.Attribute) []any {
	return []any{
		PositionComponent{X: x, Y: y},
		RenderComponent{Character: SummonPickup},
		ColorComponent{color: *color.New(colorAttr), colorAttr: colorAttr},
		SummonPickupComponent{color: *color.New(colorAttr), colorAttr: colorAttr},
	}
}
