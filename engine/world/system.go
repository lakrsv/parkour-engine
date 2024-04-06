package world

type System interface {
}

type InitializeSystem interface {
	System
	Initialize(world *World) error
}

type UpdateSystem interface {
	System
	Update(world *World) error
}
