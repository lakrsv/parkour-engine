package world

import (
	"context"
	"io"
	"log"
	"runtime/debug"
	"sync"
)

type World struct {
	cancel  context.CancelFunc
	threads sync.WaitGroup
	systems map[int][]System
	//initializeSystems []InitializeSystem
	//updateSystems     []UpdateSystem
}

func NewWorld() *World {
	return &World{systems: map[int][]System{}}
}

func (world *World) AddSystems(systems ...System) *World {
	for _, system := range systems {
		switch system.(type) {
		case InitializeSystem:
			world.systems[INITIALIZE_SYSTEM] = append(world.systems[INITIALIZE_SYSTEM], system)
		case UpdateSystem:
			world.systems[UPDATE_SYSTEM] = append(world.systems[UPDATE_SYSTEM], system)
		}
	}
	world.threads.Add(len(systems))
	return world
}

func (world *World) Simulate(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	world.cancel = cancel

	// TODO: How to wait before going to update here.. sync?
	world.initialize()

	for {
		// TODO: How to ensure all systems have updated once before going on
		world.update()
	}
}

func (world *World) Close() error {
	world.threads.Add(1)
	if world.cancel != nil {
		world.cancel()
	}

	for _, system := range world.systems[INITIALIZE_SYSTEM] {
		if closer, ok := system.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
	}

	for _, system := range world.systems[UPDATE_SYSTEM] {
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
	for _, system := range world.systems[INITIALIZE_SYSTEM] {
		initialize := func() {
			defer handlePanic()
			if err := system.(InitializeSystem).Initialize(world); err != nil {
				// TODO: Handle error
				panic(err)
			}
		}
		go initialize()
	}
}

func (world *World) update() {
	for _, system := range world.systems[UPDATE_SYSTEM] {
		update := func() {
			defer handlePanic()
			if err := system.(UpdateSystem).Update(world); err != nil {
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
