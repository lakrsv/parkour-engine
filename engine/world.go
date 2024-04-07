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

// World TODO: Add Time for rate limiting goroutines (update?)
type World struct {
	cancel     context.CancelFunc
	threads    sync.WaitGroup
	systems    map[SystemType][]System
	components ComponentStorage
	Time       *Time
}

func NewWorld() *World {
	return &World{
		systems:    map[SystemType][]System{},
		Time:       newTime(),
		components: NewComponentStorage()}
}

func (world *World) AddSystems(systems ...System) *World {
	for _, system := range systems {
		switch system.(type) {
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
	return world.components.createEntity(components...)
}

func (world *World) GetGroup(m Matcher) *Group {
	return newGroup(m, &world.components)
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
	if world.cancel != nil {
		world.cancel()
	}

	for _, system := range world.systems[Initialize] {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
	}

	for _, system := range world.systems[Update] {
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
