package main

import (
	"embed"
	"golang.org/x/term"
	"os"
)

//go:embed assets/*
var content embed.FS

func main() {
	InitAudio()
	go PlayBackgroundMusic()
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	Run(0)
}
