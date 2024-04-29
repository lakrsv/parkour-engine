package engine

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type World struct {
	wg         sync.WaitGroup
	systems    map[SystemType][]System
	components *ComponentStorage
	groups     map[Matcher]*Group
	running    bool
	Time       *Time
}

func NewWorld() *World {
	return &World{
		systems:    map[SystemType][]System{},
		Time:       newTime(time.Second/60, time.Second/60),
		components: NewComponentStorage(),
		groups:     make(map[Matcher]*Group),
		running:    true,
	}
}

func (world *World) AddSystems(systems ...System) *World {
	for _, system := range systems {
		switch system.(type) {
		case InitializeUpdateSystem:
			world.systems[Initialize] = append(world.systems[Initialize], system)
			world.systems[Update] = append(world.systems[Update], system)
		case InitializeSystem:
			world.systems[Initialize] = append(world.systems[Initialize], system)
		case UpdateSystem:
			world.systems[Update] = append(world.systems[Update], system)
		}
	}
	return world
}

func (world *World) CreateEntity(components ...any) uint32 {
	entity := world.components.createEntity(components...)
	var wg sync.WaitGroup
	wg.Add(len(world.groups))
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components, &wg)
	}
	wg.Wait()
	return entity
}

func (world *World) DeleteEntity(entity uint32) {
	world.components.deleteEntity(entity)
	var wg sync.WaitGroup
	wg.Add(len(world.groups))
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components, &wg)
	}
	wg.Wait()
}

func (world *World) GetGroup(m Matcher) *Group {
	if val, ok := world.groups[m]; ok {
		return val
	}
	world.groups[m] = newGroup(m, world.components)
	return world.groups[m]
}

func (world *World) GetUniqueComponent(t reflect.Type) any {
	if !world.components.hasComponent(t) {
		slog.Error("Component not found in component storage")
		return nil
	}
	set := world.components.getComponentSet(t)
	if set.components.Len() != 1 {
		slog.Error(fmt.Sprintf("Expected 1 entity in component set, found %d", set.components.Len()))
		return nil
	}
	_, component, _ := set.components.Iterator().Next()
	return component
}

func (world *World) ReplaceUniqueComponent(component any) {
	t := reflect.TypeOf(component)
	if !world.components.hasComponent(t) {
		slog.Error("Component not found in component storage")
		return
	}
	set := world.components.getComponentSet(t)
	if set.components.Len() != 1 {
		slog.Error(fmt.Sprintf("Expected 1 entity in component set, found %d", set.components.Len()))
		return
	}
	id, _, _ := set.components.Iterator().Next()
	world.ReplaceComponent(id, component)
}

func (world *World) GetEntityComponent(entity uint32, t reflect.Type) (any, bool) {
	if !world.components.hasComponent(t) {
		return nil, false
	}
	set := world.components.getComponentSet(t)
	if set.components.Contains(entity) {
		return set.getComponent(entity), true
	}
	return nil, false
}

func (world *World) ReplaceComponent(entity uint32, component any) {
	if !world.components.hasComponent(reflect.TypeOf(component)) {
		world.AddComponent(entity, component)
		return
	}
	set := world.components.getComponentSet(reflect.TypeOf(component))
	if !set.components.Contains(entity) {
		world.AddComponent(entity, component)
		return
	}
	set.replaceComponent(entity, component)
	var wg sync.WaitGroup
	wg.Add(len(world.groups))
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components, &wg)
	}
	wg.Wait()
}

func (world *World) AddComponent(entity uint32, component any) {
	if !world.components.hasComponent(reflect.TypeOf(component)) {
		world.components.registerComponent(reflect.TypeOf(component))
	}
	set := world.components.getComponentSet(reflect.TypeOf(component))
	if set.components.Contains(entity) {
		slog.Error("Entity already registered in component storage", "stack", getStack())
		return
	}
	set.addComponent(entity, component)
	var wg sync.WaitGroup
	wg.Add(len(world.groups))
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components, &wg)
	}
	wg.Wait()
}

func (world *World) RemoveComponent(entity uint32, t reflect.Type) {
	if !world.components.hasComponent(t) {
		return
	}
	set := world.components.getComponentSet(t)
	if !set.components.Contains(entity) {
		return
	}
	set.removeEntity(entity)
	var wg sync.WaitGroup
	wg.Add(len(world.groups))
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components, &wg)
	}
	wg.Wait()
}

func (world *World) Simulate() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	_, err = window.GetSurface()
	if err != nil {
		panic(err)
	}
	world.initialize()
	world.CreateEntity(InputComponent{KeyState: make(map[sdl.Keycode]bool)})
	for world.running {
		input := make(map[sdl.Keycode]bool)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			key, state := world.handleEvent(event)
			if key == sdl.K_UNKNOWN {
				continue
			}
			if state == sdl.RELEASED {
				input[key] = false
			} else if state == sdl.PRESSED {
				input[key] = true
			}
		}
		world.ReplaceUniqueComponent(InputComponent{KeyState: input})
		loopTime := world.loop()
		window.UpdateSurface()
		delay := uint32(world.Time.Timestep.Milliseconds()) - loopTime
		sdl.Delay(delay)
	}
	return nil
}

func (world *World) handleEvent(event sdl.Event) (sdl.Keycode, uint8) {
	switch t := event.(type) {
	case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
		println("Quitting..")
		world.running = false
		break
	case *sdl.KeyboardEvent:
		return t.Keysym.Sym, t.State
	}
	return sdl.K_UNKNOWN, 0
}

func (world *World) loop() uint32 {
	startTime := sdl.GetTicks()

	world.update()
	world.Time.update()

	return sdl.GetTicks() - startTime
}

func (world *World) Close() error {
	if !world.running {
		return nil
	}
	world.running = false

	systems := map[System]bool{}
	for _, s := range world.systems {
		for _, system := range s {
			systems[system] = true
		}
	}
	for system := range systems {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				slog.Error(
					"Failed closing system %s",
					"stack", getStack(),
				)
			}
		}
	}
	return nil
}

func (world *World) initialize() {
	for _, system := range world.systems[Initialize] {
		initialize := func() {
			defer handlePanic()
			if err := system.(InitializeSystem).Initialize(world); err != nil {
				slog.Error(
					"Failed initializing system",
					"stack", getStack(),
				)
			}
		}
		initialize()
	}
}

func (world *World) update() {
	for _, system := range world.systems[Update] {
		update := func() {
			defer handlePanic()
			if err := system.(UpdateSystem).Update(world); err != nil {
				slog.Error(
					"Failed updating system",
					"stack", debug.Stack(),
				)
			}
		}
		update()
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		slog.Error(
			"Panic",
			"recover", r,
			"stack", getStack(),
		)
	}
}

func getStack() []byte {
	buf := make([]byte, 1<<20)
	len := runtime.Stack(buf, true)
	return buf[:len]
}
