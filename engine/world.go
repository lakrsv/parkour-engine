package engine

import (
	"context"
	"golang.org/x/time/rate"
	"io"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"
	"time"
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
		Time:       newTime(0, time.Second/60),
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
				slog.Error("Failed closing system %s", slog.Any("stack", debug.Stack()))
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
				slog.Error("Failed initializing system", slog.Any("stack", debug.Stack()))
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
				slog.Error("Failed updating system", slog.Any("stack", debug.Stack()))
			}
		}
		update()
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		slog.Error("Panic", slog.Any("recover", r), slog.Any("stack", debug.Stack()))
	}
}
