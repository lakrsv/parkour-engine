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
	entityIndex   int
	entities      *intsets.Sparse
	componentSets []ComponentSet[any]
}

func NewComponentStorage() *ComponentStorage {
	return &ComponentStorage{
		registry:      map[reflect.Type]int{},
		entities:      &intsets.Sparse{},
		componentSets: []ComponentSet[any]{},
	}
}

func (storage *ComponentStorage) registerComponent(t reflect.Type) {
	idx := len(storage.registry)
	storage.registry[t] = idx
	storage.componentSets = append(storage.componentSets, ComponentSet[any]{entities: &intsets.Sparse{}, components: make([]any, MaxEntitiesPerComponent)})
}

func (storage *ComponentStorage) getComponentId(t reflect.Type) int {
	return storage.registry[t]
}

func (storage *ComponentStorage) getComponentSet(t reflect.Type) *ComponentSet[any] {
	return &storage.componentSets[storage.getComponentId(t)]
}

func (storage *ComponentStorage) createEntity(components ...any) int {
	entity := storage.entityIndex
	storage.entityIndex++

	storage.entities.Insert(entity)

	for _, component := range components {
		set := storage.getComponentSet(reflect.TypeOf(component))
		set.addEntityComponent(entity, component)
	}
	return entity
}

func (storage *ComponentStorage) deleteEntity(entity int) {
	storage.entities.Remove(entity)
	for _, set := range storage.componentSets {
		set.entities.Remove(entity)
		set.components[entity] = nil
	}
}

type ComponentSet[T any] struct {
	entities   *intsets.Sparse
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
