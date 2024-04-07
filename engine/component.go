package engine

import (
	"golang.org/x/tools/container/intsets"
	"reflect"
)

const (
	MaxEntitiesPerComponent = 8192
)

type ComponentStorage struct {
	registry      map[reflect.Type]int
	componentSets []ComponentSet[any]
}

func NewComponentStorage() ComponentStorage {
	return ComponentStorage{registry: map[reflect.Type]int{}, componentSets: []ComponentSet[any]{}}
}

func (storage *ComponentStorage) registerComponent(t reflect.Type) {
	idx := len(storage.registry)
	storage.registry[t] = idx
	storage.componentSets = append(storage.componentSets, ComponentSet[any]{entities: intsets.Sparse{}, components: make([]any, MaxEntitiesPerComponent)})
}

func (storage *ComponentStorage) getComponentId(t reflect.Type) int {
	return storage.registry[t]
}

func (storage *ComponentStorage) getComponentSet(t reflect.Type) *ComponentSet[any] {
	return &storage.componentSets[storage.registry[t]]
}

type ComponentSet[T any] struct {
	entities   intsets.Sparse
	components []T
}

func (set *ComponentSet[T]) addEntityComponent(entity int, component T) {
	if set.entities.Has(entity) {
		// TODO: Handle error
		panic("Entity already exists")
	}
	set.entities.Insert(entity)
	set.components[entity] = component
}

func (set *ComponentSet[T]) hasEntity(entity int) bool {
	return set.entities.Has(entity)
}
