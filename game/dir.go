package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

func getBaseDirectory() string {
	var baseDirectory string
	if dir, err := os.Executable(); err != nil {
		slog.Error("Failed getting path to executable", "error", err)
		baseDirectory = "."
	} else {
		baseDirectory = filepath.Dir(dir)
	}
	baseDirectory += "/assets"
	return baseDirectory
}
