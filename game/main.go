package main

import (
	"context"
	"github.com/lakrsv/parkour/engine"
	"golang.org/x/tools/container/intsets"
	"log"
	"reflect"
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
	w := engine.NewWorld()

	w.RegisterComponent(reflect.TypeOf(Test1Component{}))
	w.RegisterComponent(reflect.TypeOf(Test2Component{}))

	w.AddSystems(&HelloSystem{}, &UpdateSystem{}, &UpdateSystem{})

	entity := w.CreateEntity(
		Test1Component{x: 0, y: 10, z: 200},
		Test2Component{msg: "Hello World!"},
	)

	log.Printf("Created entity %v", entity)

	result := w.Query(engine.EntityQuery{WithAllComponents: []reflect.Type{
		reflect.TypeOf(Test1Component{}),
		reflect.TypeOf(Test2Component{}),
	}})

	for {
		val := result.Min()
		if val == intsets.MaxInt {
			break
		}
		log.Printf("Got entity: %v", val)
		result.Remove(val)
	}

	if err := w.Simulate(context.Background()); err != nil {
		panic(err)
	}
}

type HelloSystem struct {
}

func (helloSystem *HelloSystem) Initialize(world *engine.World) error {
	log.Println("Hello from HelloSystem!")
	return nil
}

type UpdateSystem struct {
	num int
}

func (updateSystem *UpdateSystem) Update(world *engine.World) error {
	updateSystem.num += 1
	if updateSystem.num > 10 {
		world.Close()
	}
	log.Printf("Time elapsed %v", world.Time.DeltaTime)
	return nil
}
