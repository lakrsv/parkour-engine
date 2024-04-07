package main

import (
	"context"
	"github.com/lakrsv/parkour/engine"
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
	w.RegisterComponent(reflect.TypeOf(Test3Component{}))

	w.AddSystems(&HelloSystem{}, &UpdateSystem{}, &UpdateSystem{})

	_ = w.CreateEntity(
		Test1Component{x: 0, y: 10, z: 200},
		Test2Component{msg: "Hello World!"},
		Test3Component{a: 5},
	)
	_ = w.CreateEntity(
		Test1Component{x: 10, y: 20, z: 300},
	)

	// AllOfComponent Manual Test
	log.Println("Running AllOfComponent Manual Test (Expecting [0])")
	log.Println(w.GetGroup(&engine.AllOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(Test1Component{}),
		reflect.TypeOf(Test2Component{}),
	}}).GetEntities())

	// AnyOfComponent Manual Test
	log.Println("Running AnyOfComponent Manual Test (Expecting [0 1])")
	log.Println(w.GetGroup(&engine.AnyOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(Test1Component{}),
		reflect.TypeOf(Test2Component{}),
	}}).GetEntities())

	// NoneOfComponent Manual Test
	log.Println("Running NoneOfComponent Manual Test (Expecting [1])")
	log.Println(w.GetGroup(&engine.NoneOfComponentMatcher{Components: []reflect.Type{
		reflect.TypeOf(Test2Component{}),
	}}).GetEntities())

	// AllOf Manual Test
	log.Println("Running AllOf Manual Test (Expecting [])")
	log.Println(w.GetGroup(&engine.AllOfMatcher{Matchers: []engine.Matcher{
		&engine.AllOfComponentMatcher{Components: []reflect.Type{
			reflect.TypeOf(Test1Component{}),
			reflect.TypeOf(Test2Component{}),
		}},
		&engine.NoneOfComponentMatcher{Components: []reflect.Type{
			reflect.TypeOf(Test3Component{}),
		}},
	}}).GetEntities())

	if err := w.Simulate(context.Background()); err != nil {
		panic(err)
	}
}

type HelloSystem struct {
}

func (helloSystem *HelloSystem) Initialize(_ *engine.World) error {
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
