module github.com/lakrsv/parkour-engine/game

go 1.24.1

require (
	atomicgo.dev/cursor v0.2.0
	github.com/gopxl/beep v1.4.1
	github.com/lakrsv/parkour-engine/engine v0.0.0-00010101000000-000000000000
	github.com/veandco/go-sdl2 v0.4.40
)

require (
	github.com/ebitengine/oto/v3 v3.2.0 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace github.com/lakrsv/parkour-engine/engine => ../engine
