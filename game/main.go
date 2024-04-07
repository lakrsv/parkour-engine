package main

import (
	"context"
	"github.com/lakrsv/parkour/engine/world"
	"log"
)

//import (
//	"github.com/lakrsv/parkour/engine"
//	"os"
//)
//
//var winTitle = "Go-SDL2 Texture"
//var winWidth, winHeight int32 = 800, 600
//var imageName = "./assets/test.png"
//
//func main() {
//	world := engine.NewWorld(winWidth, winHeight)
//	sprite, err := engine.NewSprite(world, imageName)
//	if err != nil {
//		panic(err)
//	}
//	spriteRenderer, err := engine.NewSpriteRenderer(sprite)
//	if err != nil {
//		panic(err)
//	}
//	world.AddRenderer(spriteRenderer)
//
//	world.Draw()
//	os.Exit(0)
//}

func main() {
	w := world.NewWorld()
	w.AddSystems(&HelloSystem{}, &UpdateSystem{}, &UpdateSystem{})

	if err := w.Simulate(context.Background()); err != nil {
		panic(err)
	}
}

type HelloSystem struct {
}

func (helloSystem *HelloSystem) Initialize(world *world.World) error {
	println("Hello world!")
	return nil
}

type UpdateSystem struct {
	num int
}

func (updateSystem *UpdateSystem) Update(world *world.World) error {
	updateSystem.num += 1
	if updateSystem.num > 10000 {
		world.Close()
	}
	log.Printf("Time elapsed %v", world.Time.DeltaTime)
	return nil
}
