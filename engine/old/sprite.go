package old

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type SpriteRenderer struct {
	sprite Sprite
}

type Sprite struct {
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func NewSprite(world *World, resourcePath string) (*Sprite, error) {
	var renderer *sdl.Renderer
	var texture *sdl.Texture
	var err error
	renderer, err = sdl.CreateRenderer(world.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return nil, err
	}
	// TODO: Defer destroy?
	// defer renderer.Destroy()

	image, err := img.Load(resourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		return nil, err
	}
	// TODO: Defer free?
	//defer image.Free()

	texture, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		return nil, err
	}
	// TODO: Defer destroy?
	// defer texture.Destroy()

	return &Sprite{renderer, texture}, nil
}

func NewSpriteRenderer(sprite *Sprite) (*SpriteRenderer, error) {
	return &SpriteRenderer{*sprite}, nil
}

func (spriteRenderer *SpriteRenderer) render(world *World) {
	var src, dst sdl.Rect
	// TODO: ??
	src = sdl.Rect{0, 0, 512, 512}
	dst = sdl.Rect{100, 50, 512, 512}

	// TODO: Don't reach into sprite.renderer?
	spriteRenderer.sprite.renderer.Clear()
	spriteRenderer.sprite.renderer.SetDrawColor(255, 0, 0, 255)
	spriteRenderer.sprite.renderer.FillRect(&sdl.Rect{0, 0, world.Width, world.Height})
	spriteRenderer.sprite.renderer.Copy(spriteRenderer.sprite.texture, &src, &dst)
	spriteRenderer.sprite.renderer.Present()
}
