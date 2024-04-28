package src

import "embed"

var content embed.FS

func SetAssets(c embed.FS) {
	content = c
}
