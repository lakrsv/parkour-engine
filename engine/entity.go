package engine

type EntityIndex uint64

type Entity struct {
	id EntityIndex
}

func newEntity(id EntityIndex) Entity {
	return Entity{id: id}
}
