package world

import (
	"context"
	"io"
	"log"
	"runtime/debug"
	"sync"
)

type World struct {
	cancel            context.CancelFunc
	threads           sync.WaitGroup
	initializeSystems []InitializeSystem
	updateSystems     []UpdateSystem
	systems           []System
}

func NewWorld() *World {
	return &World{initializeSystems: []InitializeSystem{}, updateSystems: []UpdateSystem{}}
}

func (world *World) AddSystems(systems ...System) *World {
	world.systems = append(world.systems, systems...)
	world.threads.Add(len(systems))
	return world
}

func (world *World) AddInitializeSystems(systems ...InitializeSystem) *World {
	world.initializeSystems = append(world.initializeSystems, systems...)
	world.threads.Add(len(systems))
	return world
}

func (world *World) AddUpdateSystems(systems ...UpdateSystem) *World {
	world.updateSystems = append(world.updateSystems, systems...)
	world.threads.Add(len(systems))
	return world
}

func (world *World) Simulate(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	world.cancel = cancel

	world.initialize()
	for {
		world.update()
	}
}

func (world *World) Close() error {
	world.threads.Add(1)
	if world.cancel != nil {
		world.cancel()
	}

	for _, system := range world.initializeSystems {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
	}

	for _, system := range world.updateSystems {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
	}

	world.threads.Done()
	world.threads.Wait()
	return nil
}

func (world *World) initialize() {
	for _, system := range world.initializeSystems {
		initialize := func() {
			defer handlePanic()
			if err := system.Initialize(world); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
		go initialize()
	}
}

func (world *World) update() {
	for _, system := range world.updateSystems {
		update := func() {
			defer handlePanic()
			if err := system.Update(world); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
		go update()
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("panic: %s \n %s", r, debug.Stack())
	}
}
