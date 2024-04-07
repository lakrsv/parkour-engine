package old

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type World struct {
	Width, Height int32
	Window        *sdl.Window
	renderers     []Renderer
}

func NewWorld(width, height int32) *World {
	// TODO: Parameterise
	window, err := sdl.CreateWindow("Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		panic(err)
	}
	return &World{width, height, window, []Renderer{}}
}

func (world *World) Draw() {
	for _, renderer := range world.renderers {
		renderer.render(world)
	}
	sdl.PollEvent()
	sdl.Delay(2000)
}

func (world *World) AddRenderer(renderer Renderer) {
	world.renderers = append(world.renderers, renderer)
}
