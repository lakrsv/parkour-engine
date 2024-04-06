package main

import (
	"github.com/lakrsv/parkour/engine"
	"os"
)

var winTitle = "Go-SDL2 Texture"
var winWidth, winHeight int32 = 800, 600
var imageName = "./assets/test.png"

func main() {
	world := engine.NewWorld(winWidth, winHeight)
	sprite, err := engine.NewSprite(world, imageName)
	if err != nil {
		panic(err)
	}
	spriteRenderer, err := engine.NewSpriteRenderer(sprite)
	if err != nil {
		panic(err)
	}
	world.AddRenderer(spriteRenderer)

	world.Draw()
	os.Exit(0)
}
