package main

import (
	"github.com/lakrsv/parkour/engine/render"
	"os"
)

var winTitle = "Go-SDL2 Texture"
var winWidth, winHeight int32 = 800, 600
var imageName = "./assets/test.png"

func main() {
	world := render.NewWorld(winWidth, winHeight)
	sprite, err := render.NewSprite(world, imageName)
	if err != nil {
		panic(err)
	}
	spriteRenderer, err := render.NewSpriteRenderer(sprite)
	if err != nil {
		panic(err)
	}
	world.AddRenderer(spriteRenderer)

	world.Draw()
	os.Exit(0)
}
