package engine

import (
	"log/slog"
	"reflect"
	"runtime/debug"
)

const (
	MaxEntities = 8192
)

type ComponentStorage struct {
	registry      map[reflect.Type]int
	entityIndex   uint32
	entities      *SparseSet[any]
	componentSets []ComponentSet[any]
}

func NewComponentStorage() *ComponentStorage {
	return &ComponentStorage{
		registry:      map[reflect.Type]int{},
		entities:      NewSparseSet[any](MaxEntities),
		componentSets: []ComponentSet[any]{},
	}
}

func (storage *ComponentStorage) hasComponent(t reflect.Type) bool {
	_, ok := storage.registry[t]
	return ok
}

func (storage *ComponentStorage) registerComponent(t reflect.Type) {
	idx := len(storage.registry)
	storage.registry[t] = idx
	storage.componentSets = append(storage.componentSets, ComponentSet[any]{components: NewSparseSet[any](MaxEntities)})
}

func (storage *ComponentStorage) getComponentId(t reflect.Type) int {
	if !storage.hasComponent(t) {
		storage.registerComponent(t)
	}
	return storage.registry[t]
}

func (storage *ComponentStorage) getComponentSet(t reflect.Type) *ComponentSet[any] {
	return &storage.componentSets[storage.getComponentId(t)]
}

func (storage *ComponentStorage) createEntity(components ...any) uint32 {
	entity := storage.entityIndex
	storage.entityIndex++

	storage.entities.Insert(entity, entity)

	for _, component := range components {
		componentType := reflect.TypeOf(component)
		if _, ok := storage.registry[componentType]; !ok {
			storage.registerComponent(componentType)
		}
		set := storage.getComponentSet(componentType)
		set.addComponent(entity, component)
	}
	return entity
}

func (storage *ComponentStorage) deleteEntity(entity uint32) {
	storage.entities.Remove(entity)
	for _, set := range storage.componentSets {
		set.components.Remove(entity)
	}
}

type ComponentSet[T comparable] struct {
	components *SparseSet[T]
}

func (set *ComponentSet[T]) replaceComponent(entity uint32, component T) {
	if !set.components.Contains(entity) {
		slog.Error(
			"Entity not in componentSet",
			"entity", entity,
			"stack", debug.Stack(),
		)
		return
	}
	set.components.Remove(entity)
	set.components.Insert(entity, component)
}

func (set *ComponentSet[T]) addComponent(entity uint32, component T) {
	if set.components.Contains(entity) {
		slog.Error(
			"Entity already in componentSet",
			"entity", entity,
			"stack", debug.Stack(),
		)
		return
	}
	set.components.Insert(entity, component)
}

func (set *ComponentSet[T]) removeEntity(entity uint32) {
	if !set.components.Contains(entity) {
		slog.Error(
			"Entity not in componentSet",
			"entity", entity,
			"stack", debug.Stack(),
		)
		panic("Entity not in componentSet")
	}
	set.components.Remove(entity)
}

func (set *ComponentSet[T]) getComponent(entity uint32) T {
	if !set.components.Contains(entity) {
		slog.Error(
			"Entity not in componentSet",
			"entity", entity,
			"stack", debug.Stack(),
		)
		panic("Entity not in componentSet")
	}
	component, _ := set.components.Get(entity)
	return component
}
