package world

type System interface {
}

type SystemType int

const (
	INITIALIZE_SYSTEM SystemType = iota
	UPDATE_SYSTEM                = iota
)

type InitializeSystem interface {
	System
	Initialize(world *World) error
}

type UpdateSystem interface {
	System
	Update(world *World) error
}
