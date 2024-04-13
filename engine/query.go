package engine

import (
	"reflect"
	"sync"
)

type Matcher interface {
	match(storage *ComponentStorage) *SparseSet[uint32]
	matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32]
}

type AllOfComponentMatcher struct {
	Components []reflect.Type
}

type AnyOfComponentMatcher struct {
	Components []reflect.Type
}

type NoneOfComponentMatcher struct {
	Components []reflect.Type
}

type AllOfMatcher struct {
	Matchers []Matcher
}

type AnyOfMatcher struct {
	Matchers []Matcher
}

type NoneOfMatcher struct {
	Matchers []Matcher
}

func (m *AllOfComponentMatcher) match(storage *ComponentStorage) *SparseSet[uint32] {
	result := NewSparseSet[uint32](0)
	for i, t := range m.Components {
		set := storage.getComponentSet(t)
		if result.IsEmpty() && i == 0 {
			result = result.UnionId(set.components.CopyId())
		} else {
			result = result.IntersectId(set.components.CopyId())
		}
	}
	return result
}

func (m *AnyOfComponentMatcher) match(storage *ComponentStorage) *SparseSet[uint32] {
	result := NewSparseSet[uint32](0)
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result = result.UnionId(set.components.CopyId())
	}
	return result
}

func (m *NoneOfComponentMatcher) match(storage *ComponentStorage) *SparseSet[uint32] {
	result := NewSparseSet[uint32](0)
	result = result.UnionId(storage.entities.CopyId())
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result = result.DifferenceId(set.components.CopyId())
	}
	return result
}

func (m *AllOfMatcher) match(storage *ComponentStorage) *SparseSet[uint32] {
	result := NewSparseSet[uint32](0)
	for i, mx := range m.Matchers {
		if i == 0 {
			result = mx.match(storage)
		} else {
			match := mx.match(storage)
			result = result.IntersectId(match)
		}
	}
	return result
}

func (m *AnyOfMatcher) match(storage *ComponentStorage) *SparseSet[uint32] {
	result := NewSparseSet[uint32](0)
	for _, mx := range m.Matchers {
		result = result.UnionId(mx.match(storage)).CopyId()
	}
	return result
}

func (m *NoneOfMatcher) match(storage *ComponentStorage) *SparseSet[uint32] {
	result := NewSparseSet[uint32](0)
	result = result.UnionId(storage.entities.CopyId())
	for _, mx := range m.Matchers {
		result = result.DifferenceId(mx.match(storage)).CopyId()
	}
	return result
}

func (m *AllOfComponentMatcher) matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32] {
	result := NewSparseSet[uint32](entity + 1)
	if len(m.Components) == 0 {
		return result
	}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if !set.components.Contains(entity) {
			result.Remove(entity)
			return result
		}
	}
	result.Insert(entity, entity)
	return result
}

func (m *AnyOfComponentMatcher) matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32] {
	result := NewSparseSet[uint32](entity + 1)
	if len(m.Components) == 0 {
		return result
	}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if set.components.Contains(entity) {
			result.Insert(entity, entity)
			return result
		}
	}
	result.Remove(entity)
	return result
}

func (m *NoneOfComponentMatcher) matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32] {
	result := NewSparseSet[uint32](entity + 1)
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if set.components.Contains(entity) {
			result.Remove(entity)
			return result
		}
	}
	result.Insert(entity, entity)
	return result
}

func (m *AllOfMatcher) matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32] {
	result := NewSparseSet[uint32](entity + 1)
	for i, mx := range m.Matchers {
		if i == 0 {
			result = mx.matchOne(storage, entity)
		} else {
			result = result.IntersectId(mx.matchOne(storage, entity))
		}
	}
	return result
}

func (m *AnyOfMatcher) matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32] {
	result := NewSparseSet[uint32](entity + 1)
	for _, mx := range m.Matchers {
		result = result.UnionId(mx.matchOne(storage, entity))
	}
	return result
}

func (m *NoneOfMatcher) matchOne(storage *ComponentStorage, entity uint32) *SparseSet[uint32] {
	result := NewSparseSet[uint32](entity + 1)
	result.Insert(entity, entity)
	for _, mx := range m.Matchers {
		result = result.DifferenceId(mx.matchOne(storage, entity))
	}
	return result
}

type Group struct {
	matcher       Matcher
	result        *SparseSet[uint32]
	EntityAdded   chan uint32
	EntityRemoved chan uint32
}

func newGroup(matcher Matcher, storage *ComponentStorage) *Group {
	return &Group{
		matcher:       matcher,
		result:        matcher.match(storage),
		EntityAdded:   make(chan uint32),
		EntityRemoved: make(chan uint32),
	}
}

func (g *Group) GetEntities() []uint32 {
	entities := make([]uint32, g.result.Len())
	iterator := g.result.Iterator()
	idx := 0
	for {
		id, _, ok := iterator.Next()
		if !ok {
			break
		}
		entities[idx] = id
		idx++
	}
	return entities
}

func (g *Group) EvaluateEntity(entity uint32, storage *ComponentStorage, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		if !storage.entities.Contains(entity) {
			if g.result.Contains(entity) {
				g.result.Remove(entity)
				select {
				case g.EntityRemoved <- entity:
				default:
				}
			}
			return
		}
		result := g.matcher.matchOne(storage, entity)

		if result.Contains(entity) && g.result.Contains(entity) {
			return
		}
		if !result.Contains(entity) && !g.result.Contains(entity) {
			return
		}

		if result.Contains(entity) {
			g.result.Insert(entity, entity)
			select {
			case g.EntityAdded <- entity:
			default:
			}
		} else {
			g.result.Remove(entity)
			select {
			case g.EntityRemoved <- entity:
			default:
			}
		}
	}()
}
