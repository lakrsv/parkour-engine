package engine

type Color32 struct {
	R, G, B, A uint8
}

var (
	Black = Color32{0, 0, 0, 255}
	White = Color32{255, 255, 255, 255}
	Red   = Color32{255, 0, 0, 255}
	Green = Color32{0, 255, 0, 255}
	Blue  = Color32{0, 0, 255, 255}
)
