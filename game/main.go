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

	_ = w.CreateEntity(
		Test1Component{x: 0, y: 10, z: 200},
		Test2Component{msg: "Hello World!"},
	)
	_ = w.CreateEntity(
		Test1Component{x: 10, y: 20, z: 300},
	)

	// AllOf Manual Test
	log.Println("Running AllOf Manual Test (Expecting id 0)")
	printResultTest(w.GetEntities(&engine.AllOfMatcher{Components: []reflect.Type{
		reflect.TypeOf(Test1Component{}),
		reflect.TypeOf(Test2Component{}),
	}}))

	// AnyOf Manual Test
	log.Println("Running AnyOf Manual Test (Expecting id 0 and 1)")
	printResultTest(w.GetEntities(&engine.AnyOfMatcher{Components: []reflect.Type{
		reflect.TypeOf(Test1Component{}),
		reflect.TypeOf(Test2Component{}),
	}}))

	// NoneOf Manual Test
	log.Println("Running NoneOf Manual Test (Expecting id 1)")
	printResultTest(w.GetEntities(&engine.NoneOfMatcher{Components: []reflect.Type{
		reflect.TypeOf(Test2Component{}),
	}}))

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
		if err := world.Close(); err != nil {
			panic(err)
		}
	}
	log.Printf("Time elapsed %v", world.Time.DeltaTime)
	return nil
}

func printResultTest(result *intsets.Sparse) {
	for {
		val := result.Min()
		if val == intsets.MaxInt {
			break
		}
		log.Printf("Got entity with id: %v", val)
		result.Remove(val)
	}
}
