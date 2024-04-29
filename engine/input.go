package engine

import "github.com/veandco/go-sdl2/sdl"

type InputComponent struct {
	KeyState map[sdl.Keycode]bool
}

func (c *InputComponent) KeyPressed(key sdl.Keycode) bool {
	state, ok := c.KeyState[key]
	if ok {
		return state
	}
	return false
}

func (c *InputComponent) KeyReleased(key sdl.Keycode) bool {
	state, ok := c.KeyState[key]
	if ok {
		return !state
	}
	return false
}
