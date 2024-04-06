package render

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColorBlack(t *testing.T) {
	assert.Equal(t, Black, Color32{0, 0, 0, 255})
}

func TestColorWhite(t *testing.T) {
	assert.Equal(t, White, Color32{255, 255, 255, 255})
}

func TestColorRed(t *testing.T) {
	assert.Equal(t, Red, Color32{255, 0, 0, 255})
}

func TestColorGreen(t *testing.T) {
	assert.Equal(t, Green, Color32{0, 255, 0, 255})
}

func TestColorBlue(t *testing.T) {
	assert.Equal(t, Blue, Color32{0, 0, 255, 255})
}
