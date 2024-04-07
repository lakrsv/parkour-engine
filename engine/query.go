package engine

import (
	"golang.org/x/tools/container/intsets"
	"reflect"
)

type Query interface {
	execute(world *World) *intsets.Sparse
}
type EntityQuery struct {
	WithAllComponents []reflect.Type
}

func (q EntityQuery) execute(world *World) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, t := range q.WithAllComponents {
		set := world.components.getComponentSet(t)
		l := set.entities.Len()
		if result.IsEmpty() {
			result.Copy(&set.entities)
		} else {
			result.IntersectionWith(&set.entities)
		}
	}
	return result
}
