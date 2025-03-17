package main

import (
	"embed"
)

//go:embed assets/*
var content embed.FS

func main() {
	InitAudio()
	go PlayBackgroundMusic()
	Run(0)
}
