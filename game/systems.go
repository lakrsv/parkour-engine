package main

import (
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"time"

	"atomicgo.dev/cursor"
	"github.com/lakrsv/parkour-engine/engine"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type RenderSystem struct {
	palette RunePalette
}

func (s *RenderSystem) Initialize(w *engine.World) error {
	cursor.Hide()
	for range 64 {
		fmt.Println()
	}
	return nil
}

func (s *RenderSystem) Update(w *engine.World) error {
	grid := reflect.ValueOf(w.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)
	level := reflect.ValueOf(w.GetUniqueComponent(reflect.TypeOf(LevelComponent{}))).Interface().(LevelComponent)

	var font *ttf.Font
	var surface *sdl.Surface

	surface, _ = w.Window.GetSurface()

	// Load the font for our text
	font, _ = ttf.OpenFont("./assets/fonts/consolas.ttf", 16)
	defer font.Close()

	// Render header text
	for i, headerLine := range level.Header {
		text, _ := font.RenderUTF8Blended(headerLine, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		defer text.Free()
		_ = text.Blit(nil, surface, &sdl.Rect{X: 0, Y: int32(i * 16), W: 0, H: 0})
	}

	// Calculate header height for offset
	headerHeight := len(level.Header) * 16

	// Track current x position for each row
	currentX := make([]int32, grid.Height)
	for y := range grid.Height {
		currentX[y] = 0
		for x := range grid.Width {
			entity := grid.EffectEntities[grid.GetCell(x, y)]
			if _, ok := w.GetEntityComponent(entity, reflect.TypeOf(RenderComponent{})); !ok {
				entity = math.MaxUint32
			}
			if entity == math.MaxUint32 {
				entity = grid.ForegroundEntities[grid.GetCell(x, y)]
				if _, ok := w.GetEntityComponent(entity, reflect.TypeOf(RenderComponent{})); !ok {
					entity = math.MaxUint32
				}
			}
			if entity == math.MaxUint32 {
				entity = grid.BackgroundEntities[grid.GetCell(x, y)]
			}
			if component, ok := w.GetEntityComponent(entity, reflect.TypeOf(RenderComponent{})); ok {
				render := reflect.ValueOf(component).Interface().(RenderComponent)
				var textColor sdl.Color
				if colorComponent, ok := w.GetEntityComponent(entity, reflect.TypeOf(ColorComponent{})); ok {
					c := reflect.ValueOf(colorComponent).Interface().(ColorComponent).color
					textColor = sdl.Color{R: c.R, G: c.G, B: c.B, A: 255}
				} else {
					c := s.palette.GetColor(render.Character)
					textColor = sdl.Color{R: c.R, G: c.G, B: c.B, A: 255}
				}
				text, _ := font.RenderUTF8Blended(string(render.Character), textColor)
				defer text.Free()
				_ = text.Blit(nil, surface, &sdl.Rect{X: currentX[y], Y: int32(y*16 + headerHeight), W: 0, H: 0})
				currentX[y] += text.W
			}
		}
	}

	// Render UI text
	uiTexts := []string{
		"WASD = Move",
		"E = Summon",
		"R = Restart",
		"Q = Quit",
	}

	// Use a new variable for UI text x-position
	for i, txt := range uiTexts {
		text, _ := font.RenderUTF8Blended(txt, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		defer text.Free()
		_ = text.Blit(nil, surface, &sdl.Rect{X: 0, Y: int32(grid.Height*16 + headerHeight + i*16), W: 0, H: 0})
	}

	return nil
}

type PlayerInputSystem struct {
	group *engine.Group
}

func (p *PlayerInputSystem) Initialize(world *engine.World) error {
	p.group = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(PlayerInputComponent{}),
		reflect.TypeOf(MoveComponent{}),
	}})
	return nil
}

func (p *PlayerInputSystem) Update(world *engine.World) error {
	input := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(engine.InputComponent{}))).Interface().(engine.InputComponent)
	if input.KeyPressed(sdl.K_q) {
		if err := world.Close(); err != nil {
			panic(err)
		}
		return nil
	}
	if input.KeyPressed(sdl.K_r) {
		currentLevel := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(LevelComponent{}))).Interface().(LevelComponent).Level
		if err := world.Close(); err != nil {
			panic(err)
		}
		Run(currentLevel)
		return nil
	}

	// Movement
	var x, y int
	if input.KeyPressed(sdl.K_w) {
		y = -1
	} else if input.KeyPressed(sdl.K_s) {
		y = 1
	}
	if input.KeyPressed(sdl.K_d) {
		x = 1
	} else if input.KeyPressed(sdl.K_a) {
		x = -1
	}
	if x != 0 || y != 0 {
		for _, entity := range p.group.GetEntities() {
			world.ReplaceComponent(entity, MoveComponent{x, y})
		}
	}

	if input.KeyPressed(sdl.K_e) {
		for _, entity := range p.group.GetEntities() {
			if _, ok := world.GetEntityComponent(entity, reflect.TypeOf(CreateSummonComponent{})); !ok {
				world.AddComponent(entity, CreateSummonComponent{})
			}
		}
	}

	return nil
}

type MoveSystem struct {
	group *engine.Group
}

func (m *MoveSystem) Initialize(world *engine.World) error {
	m.group = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(MoveComponent{}),
		reflect.TypeOf(PositionComponent{}),
	}})
	return nil
}

func (m *MoveSystem) Update(world *engine.World) error {
	grid := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)
	for _, entity := range m.group.GetEntities() {
		if moveComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(MoveComponent{})); ok {
			move := reflect.ValueOf(moveComponent).Interface().(MoveComponent)
			if move.X == 0 && move.Y == 0 {
				continue
			}

			if positionComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(PositionComponent{})); ok {
				position := reflect.ValueOf(positionComponent).Interface().(PositionComponent)
				newPosition := PositionComponent{position.X + move.X, position.Y + move.Y}

				newCell := grid.GetCell(newPosition.X, newPosition.Y)

				newCellEntity := grid.BackgroundEntities[newCell]
				if _, ok := world.GetEntityComponent(newCellEntity, reflect.TypeOf(ObstacleComponent{})); ok {
					// Can not walk
					world.ReplaceComponent(entity, MoveComponent{0, 0})
					if _, ok := world.GetEntityComponent(entity, reflect.TypeOf(FacingComponent{})); ok {
						world.ReplaceComponent(entity, FacingComponent{0, 0})
					}
					continue
				}
				// Something is already there in the foreground
				if grid.ForegroundEntities[newCell] != math.MaxUint32 {
					world.ReplaceComponent(entity, MoveComponent{0, 0})
					if _, ok := world.GetEntityComponent(entity, reflect.TypeOf(FacingComponent{})); ok {
						world.ReplaceComponent(entity, FacingComponent(move))
					}
					continue
				}

				oldCell := grid.GetCell(position.X, position.Y)
				if grid.ForegroundEntities[oldCell] != entity {
					slog.Error(
						"Entity in foreground did not match expected entity!",
						"expected", entity,
						"actual", grid.ForegroundEntities[oldCell],
					)
				}
				grid.ForegroundEntities[oldCell] = math.MaxUint32
				grid.ForegroundEntities[newCell] = entity

				world.ReplaceComponent(entity, MoveComponent{0, 0})
				if _, ok := world.GetEntityComponent(entity, reflect.TypeOf(FacingComponent{})); ok {
					world.ReplaceComponent(entity, FacingComponent(move))
				}
				world.ReplaceComponent(entity, newPosition)
				playWalkSound()
			}
		}
	}
	return nil
}

type DeferDoorRenderSystem struct {
	group *engine.Group
}

func (s *DeferDoorRenderSystem) Initialize(world *engine.World) error {
	s.group = world.GetGroup(
		&engine.AllOfMatcher{Matchers: []engine.Matcher{
			&engine.AllOfComponentMatcher{Components: []reflect.Type{
				reflect.TypeOf(DeferDoorRenderComponent{}),
				reflect.TypeOf(PositionComponent{}),
			}},
			&engine.NoneOfComponentMatcher{Components: []reflect.Type{
				reflect.TypeOf(RenderComponent{}),
			}},
		}},
	)
	return nil
}

func (s *DeferDoorRenderSystem) Update(world *engine.World) error {
	grid := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)

	for _, entity := range s.group.GetEntities() {
		positionComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(PositionComponent{}))
		if !ok {
			continue
		}

		position := reflect.ValueOf(positionComponent).Interface().(PositionComponent)
		leftNeighbour := grid.BackgroundEntities[grid.GetCell(position.X-1, position.Y)]
		rightNeighbour := grid.BackgroundEntities[grid.GetCell(position.X+1, position.Y)]
		upNeighbour := grid.BackgroundEntities[grid.GetCell(position.X, position.Y-1)]
		downNeighbour := grid.BackgroundEntities[grid.GetCell(position.X, position.Y+1)]

		if _, ok := world.GetEntityComponent(upNeighbour, reflect.TypeOf(ObstacleComponent{})); ok {
			if _, ok := world.GetEntityComponent(downNeighbour, reflect.TypeOf(ObstacleComponent{})); ok {
				// Vertical Door
				world.AddComponent(entity, RenderComponent{Character: DoorVertical})
				continue
			}
		} else if _, ok := world.GetEntityComponent(leftNeighbour, reflect.TypeOf(ObstacleComponent{})); ok {
			if _, ok := world.GetEntityComponent(rightNeighbour, reflect.TypeOf(ObstacleComponent{})); ok {
				// Horizontal Door
				world.AddComponent(entity, RenderComponent{Character: DoorHorizontal})
				continue
			}
		} else {
			// Vertical Door
			world.AddComponent(entity, RenderComponent{Character: DoorVertical})
		}
	}
	return nil
}

type TriggerSystem struct {
	triggers     *engine.Group
	triggered    *engine.Group
	moving       *engine.Group
	triggeredMap map[rune]map[uint32]bool
	quit         chan func()
}

func (t *TriggerSystem) Close() error {
	close(t.quit)
	return nil
}

func (t *TriggerSystem) Initialize(world *engine.World) error {
	t.quit = make(chan func(), 1)
	t.triggers = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{reflect.TypeOf(TriggerComponent{})}})
	t.triggered = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{reflect.TypeOf(TriggeredComponent{})}})
	t.moving = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{reflect.TypeOf(PositionComponent{}), reflect.TypeOf(MoveComponent{}), reflect.TypeOf(InteractsWithTriggersComponent{})}})
	t.triggeredMap = make(map[rune]map[uint32]bool)
	for _, entity := range t.triggered.GetEntities() {
		if component, ok := world.GetEntityComponent(entity, reflect.TypeOf(TriggeredComponent{})); ok {
			triggeredComponent := reflect.ValueOf(component).Interface().(TriggeredComponent)
			if _, ok := t.triggeredMap[triggeredComponent.Symbol]; !ok {
				t.triggeredMap[triggeredComponent.Symbol] = make(map[uint32]bool, len(t.triggered.GetEntities()))
			}
			t.triggeredMap[triggeredComponent.Symbol][entity] = true
		}
	}
	go func() {
		for {
			select {
			case <-t.quit:
				return
			case id := <-t.triggered.EntityAdded:
				if component, ok := world.GetEntityComponent(id, reflect.TypeOf(TriggeredComponent{})); ok {
					triggeredComponent := reflect.ValueOf(component).Interface().(TriggeredComponent)
					if _, ok := t.triggeredMap[triggeredComponent.Symbol]; !ok {
						t.triggeredMap[triggeredComponent.Symbol] = make(map[uint32]bool, len(t.triggered.GetEntities()))
					}
					t.triggeredMap[triggeredComponent.Symbol][id] = true
				}
			case id := <-t.triggered.EntityRemoved:
				if component, ok := world.GetEntityComponent(id, reflect.TypeOf(TriggeredComponent{})); ok {
					triggeredComponent := reflect.ValueOf(component).Interface().(TriggeredComponent)
					if _, ok := t.triggeredMap[triggeredComponent.Symbol]; !ok {
						t.triggeredMap[triggeredComponent.Symbol] = make(map[uint32]bool, len(t.triggered.GetEntities()))
					}
					delete(t.triggeredMap[triggeredComponent.Symbol], id)
				}
			}
		}
	}()
	return nil
}

func (t *TriggerSystem) Update(world *engine.World) error {
	grid := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)
	for _, entity := range t.moving.GetEntities() {
		if positionComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(PositionComponent{})); ok {
			position := reflect.ValueOf(positionComponent).Interface().(PositionComponent)
			cellPos := grid.GetCell(position.X, position.Y)
			if grid.ForegroundEntities[cellPos] != entity {
				continue
			}
			backgroundEntity := grid.BackgroundEntities[cellPos]
			if triggerComponent, ok := world.GetEntityComponent(backgroundEntity, reflect.TypeOf(TriggerComponent{})); ok {

				// Check if color match
				if triggerColorComponent, ok := world.GetEntityComponent(backgroundEntity, reflect.TypeOf(ColorComponent{})); ok {
					if entityTriggerInteractColorComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(InteractsWithTriggersComponent{})); ok {
						triggerColor := reflect.ValueOf(triggerColorComponent).Interface().(ColorComponent)
						entityTriggerInteractColor := reflect.ValueOf(entityTriggerInteractColorComponent).Interface().(InteractsWithTriggersComponent)

						if triggerColor.color != entityTriggerInteractColor.color {
							continue
						}
					}
				}

				trigger := reflect.ValueOf(triggerComponent).Interface().(TriggerComponent)
				if !trigger.Triggered {
					if triggeredCol, ok := t.triggeredMap[trigger.Symbol]; ok {
						for triggeredEntity, ok := range triggeredCol {
							if !ok {
								continue
							}
							if triggeredComponent, ok := world.GetEntityComponent(triggeredEntity, reflect.TypeOf(TriggeredComponent{})); ok {
								triggered := reflect.ValueOf(triggeredComponent).Interface().(TriggeredComponent)
								triggered.Action(triggeredEntity, world)
								trigger.Triggered = true
							}
						}
					}
					world.ReplaceComponent(backgroundEntity, trigger)
					world.ReplaceComponent(backgroundEntity, RenderComponent{Character: TriggeredButton})
				}
			}
		}
	}
	return nil
}

type DirectionIndicatorSystem struct {
	facing                      *engine.Group
	directionIndicatorsByEntity map[uint32]uint32
}

func (s *DirectionIndicatorSystem) Initialize(world *engine.World) error {
	s.directionIndicatorsByEntity = make(map[uint32]uint32)
	s.facing = world.GetGroup(&engine.AllOfMatcher{Matchers: []engine.Matcher{
		&engine.AllOfComponentMatcher{Components: []reflect.Type{
			reflect.TypeOf(FacingComponent{}),
			reflect.TypeOf(PositionComponent{}),
		}},
	}})
	return nil
}

func (s *DirectionIndicatorSystem) Update(world *engine.World) error {
	for _, entity := range s.facing.GetEntities() {
		if facingComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(FacingComponent{})); ok {
			facing := reflect.ValueOf(facingComponent).Interface().(FacingComponent)
			if facing.X == 0 && facing.Y == 0 {
				if directionIndicatorEntity, ok := s.directionIndicatorsByEntity[entity]; ok {
					world.RemoveComponent(directionIndicatorEntity, reflect.TypeOf(RenderComponent{}))
				}
				continue
			}
			if positionComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(PositionComponent{})); ok {
				position := reflect.ValueOf(positionComponent).Interface().(PositionComponent)
				directionIndicator, ok := s.directionIndicatorsByEntity[entity]
				if !ok {
					directionIndicator = world.CreateEntity()
					s.directionIndicatorsByEntity[entity] = directionIndicator
				}

				grid := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)
				if directionIndicatorPositionComponent, ok := world.GetEntityComponent(directionIndicator, reflect.TypeOf(PositionComponent{})); ok {
					indicatorPosition := reflect.ValueOf(directionIndicatorPositionComponent).Interface().(PositionComponent)
					indicatorCell := grid.GetCell(indicatorPosition.X, indicatorPosition.Y)
					if grid.EffectEntities[indicatorCell] != math.MaxUint32 {
						grid.EffectEntities[indicatorCell] = math.MaxUint32
					}
				}
				newIndicatorPosition := PositionComponent{position.X + facing.X, position.Y + facing.Y}
				grid.EffectEntities[grid.GetCell(newIndicatorPosition.X, newIndicatorPosition.Y)] = directionIndicator
				world.ReplaceComponent(directionIndicator, newIndicatorPosition)

				if grid.ForegroundEntities[grid.GetCell(newIndicatorPosition.X, newIndicatorPosition.Y)] != math.MaxUint32 {
					world.RemoveComponent(directionIndicator, reflect.TypeOf(RenderComponent{}))
					continue
				}

				backgroundEntity := grid.BackgroundEntities[grid.GetCell(newIndicatorPosition.X, newIndicatorPosition.Y)]
				if _, ok := world.GetEntityComponent(backgroundEntity, reflect.TypeOf(FloorComponent{})); !ok {
					world.RemoveComponent(directionIndicator, reflect.TypeOf(RenderComponent{}))
					continue
				}

				var char rune
				if facing.Y == -1 {
					char = UpIndicator
				} else if facing.Y == 1 {
					char = DownIndicator
				} else if facing.X == 1 {
					char = RightIndicator
				} else if facing.X == -1 {
					char = LeftIndicator
				}

				if summonComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(SummonComponent{})); ok {
					summon := reflect.ValueOf(summonComponent).Interface().(SummonComponent)
					world.ReplaceComponent(directionIndicator, ColorComponent(summon))
				}

				world.ReplaceComponent(directionIndicator, RenderComponent{Character: char})
			}
		}
	}
	return nil
}

type SummonInputSystem struct {
	timePassed      time.Duration
	updateFrequency time.Duration
	group           *engine.Group
}

func (s *SummonInputSystem) Initialize(world *engine.World) error {
	s.group = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(SummonInputComponent{}),
		reflect.TypeOf(MoveComponent{}),
	}})
	s.updateFrequency = time.Second / 2
	return nil
}

func (s *SummonInputSystem) Update(world *engine.World) error {
	s.timePassed += world.Time.DeltaTime
	if s.timePassed >= s.updateFrequency {
		s.timePassed = 0
		for _, entity := range s.group.GetEntities() {
			if summonInputComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(SummonInputComponent{})); ok {
				summonInput := reflect.ValueOf(summonInputComponent).Interface().(SummonInputComponent)
				world.ReplaceComponent(entity, MoveComponent(summonInput))
			}
		}
	}
	return nil
}

type CreateSummonSystem struct {
	group *engine.Group
}

func (s *CreateSummonSystem) Initialize(world *engine.World) error {
	s.group = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(SummonComponent{}),
		reflect.TypeOf(CreateSummonComponent{}),
		reflect.TypeOf(PositionComponent{}),
		reflect.TypeOf(FacingComponent{}),
	}})
	return nil
}

func (s *CreateSummonSystem) Update(world *engine.World) error {
	for _, entity := range s.group.GetEntities() {
		if facingComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(FacingComponent{})); ok {
			facing := reflect.ValueOf(facingComponent).Interface().(FacingComponent)
			if facing.X == 0 && facing.Y == 0 {
				world.RemoveComponent(entity, reflect.TypeOf(CreateSummonComponent{}))
				continue
			}
			if positionComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(PositionComponent{})); ok {
				position := reflect.ValueOf(positionComponent).Interface().(PositionComponent)
				summonPosition := PositionComponent{position.X + facing.X, position.Y + facing.Y}

				grid := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)
				summonCell := grid.GetCell(summonPosition.X, summonPosition.Y)
				if grid.ForegroundEntities[summonCell] != math.MaxUint32 {
					world.RemoveComponent(entity, reflect.TypeOf(CreateSummonComponent{}))
					continue
				}
				if _, ok := world.GetEntityComponent(grid.BackgroundEntities[summonCell], reflect.TypeOf(ObstacleComponent{})); ok {
					world.RemoveComponent(entity, reflect.TypeOf(CreateSummonComponent{}))
					continue
				}

				if summonComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(SummonComponent{})); ok {
					summon := reflect.ValueOf(summonComponent).Interface().(SummonComponent)
					grid.ForegroundEntities[summonCell] = world.CreateEntity(
						SummonBlueprint(summonPosition.X, summonPosition.Y, facing.X, facing.Y, summon.color)...,
					)
				}
			}
		}
	}
	return nil
}

type SummonPickupSystem struct {
	group *engine.Group
}

func (s *SummonPickupSystem) Initialize(world *engine.World) error {
	s.group = world.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(SummonComponent{}),
		reflect.TypeOf(PositionComponent{}),
	}})
	return nil
}

func (s *SummonPickupSystem) Update(world *engine.World) error {
	grid := reflect.ValueOf(world.GetUniqueComponent(reflect.TypeOf(GridComponent{}))).Interface().(GridComponent)
	for _, entity := range s.group.GetEntities() {
		if positionComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(PositionComponent{})); ok {
			position := reflect.ValueOf(positionComponent).Interface().(PositionComponent)
			entityCell := grid.GetCell(position.X, position.Y)
			backgroundEntity := grid.BackgroundEntities[entityCell]
			if summonPickupComponent, ok := world.GetEntityComponent(backgroundEntity, reflect.TypeOf(SummonPickupComponent{})); ok {
				summonPickup := reflect.ValueOf(summonPickupComponent).Interface().(SummonPickupComponent)
				if summonComponent, ok := world.GetEntityComponent(entity, reflect.TypeOf(SummonComponent{})); ok {
					summon := reflect.ValueOf(summonComponent).Interface().(SummonComponent)
					if summonPickup.color != summon.color {
						playPickupSound()
						world.ReplaceComponent(entity, SummonComponent(summonPickup))
					}
				}
			}
		}
	}
	return nil
}
