package engine

type System interface {
}

//go:generate stringer -type=SystemType
type SystemType int

const (
	Initialize SystemType = iota
	Update     SystemType = iota
)

type InitializeSystem interface {
	System
	Initialize(world *World) error
}

type UpdateSystem interface {
	System
	Update(world *World) error
}

type InitializeUpdateSystem interface {
	InitializeSystem
	UpdateSystem
}
