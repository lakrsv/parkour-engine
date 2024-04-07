package world

type System interface {
}

const (
	INITIALIZE_SYSTEM = iota
	UPDATE_SYSTEM     = iota
)

type InitializeSystem interface {
	System
	Initialize(world *World) error
}

type UpdateSystem interface {
	System
	Update(world *World) error
}
