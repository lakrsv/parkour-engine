package world

type System interface {
}

// SystemType TODO: Can I force this enum onto the interface?
type SystemType int

const (
	Initialize SystemType = iota
	Update                = iota
)

type InitializeSystem interface {
	System
	Initialize(world *World) error
}

type UpdateSystem interface {
	System
	Update(world *World) error
}
