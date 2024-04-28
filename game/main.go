package main

import (
	"embed"
	"github.com/containerd/console"
	"github.com/lakrsv/parkour-engine/game/src"
)

//go:embed assets/*
var content embed.FS

func main() {
	src.SetAssets(content)
	src.InitAudio()
	go src.PlayBackgroundMusic()
	current := console.Current()
	defer current.Reset()
	if err := current.SetRaw(); err != nil {
		panic(err)
	}
	src.Run(0)
}
