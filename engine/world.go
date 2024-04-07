package engine

import (
	"context"
	"golang.org/x/time/rate"
	"io"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
)

type World struct {
	cancel     context.CancelFunc
	threads    sync.WaitGroup
	systems    map[SystemType][]System
	components *ComponentStorage
	groups     map[Matcher]*Group
	Time       *Time
}

func NewWorld() *World {
	return &World{
		systems:    map[SystemType][]System{},
		Time:       newTime(),
		components: NewComponentStorage(),
		groups:     make(map[Matcher]*Group),
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

func (world *World) RegisterComponent(t reflect.Type) {
	world.components.registerComponent(t)
}

func (world *World) CreateEntity(components ...any) int {
	entity := world.components.createEntity(components...)
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components)
	}
	return entity
}

func (world *World) DeleteEntity(entity int) {
	world.components.deleteEntity(entity)
	// TODO: Callback before or after actual deletion?
	for _, group := range world.groups {
		group.EvaluateEntity(entity, world.components)
	}
}

func (world *World) GetGroup(m Matcher) *Group {
	if val, ok := world.groups[m]; ok {
		return val
	}
	world.groups[m] = newGroup(m, world.components)
	return world.groups[m]
}

func (world *World) Simulate(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	world.cancel = cancel

	world.initialize()

	limiter := rate.NewLimiter(rate.Every(world.Time.Timestep), 1)

	world.threads.Add(1)

	go func() {
		defer world.threads.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := limiter.Wait(ctx); err != nil {
					// TODO: Handle error
					panic(err)
				}
				world.update()
				world.Time.update()
			}
		}
	}()
	world.threads.Wait()
	return nil
}

func (world *World) Close() error {
	if world.cancel == nil {
		return nil
	}
	world.cancel()
	world.cancel = nil

	systems := map[System]bool{}
	for _, s := range world.systems {
		for _, system := range s {
			systems[system] = true
		}
	}
	for system := range systems {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				// TODO: Handle error
				panic(err)
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
				// TODO: Handle error
				panic(err)
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
				// TODO: Handle error
				panic(err)
			}
		}
		update()
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("panic: %s \n %s", r, debug.Stack())
	}
}
