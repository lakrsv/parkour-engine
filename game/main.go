package main

import (
	"github.com/containerd/console"
	"github.com/lakrsv/parkour/game/src"
)

func main() {
	src.InitAudio()
	go src.PlayBackgroundMusic()
	current := console.Current()
	defer current.Reset()
	if err := current.SetRaw(); err != nil {
		panic(err)
	}
	src.Run(0)
}
