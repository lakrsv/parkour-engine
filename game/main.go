package game

import (
	"embed"
	"github.com/containerd/console"
)

//go:embed assets/*
var content embed.FS

func main() {
	InitAudio()
	go PlayBackgroundMusic()
	current := console.Current()
	defer current.Reset()
	if err := current.SetRaw(); err != nil {
		panic(err)
	}
	Run(0)
}
