package engine

import (
	"golang.org/x/tools/container/intsets"
	"reflect"
)

type Matcher interface {
	match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse
	matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse
}

type AllOfComponentMatcher struct {
	Components []reflect.Type
}

func (m *AllOfComponentMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if result == nil {
			result = &intsets.Sparse{}
			result.Copy(set.entities)
		} else {
			result.IntersectionWith(set.entities)
		}
	}
	return result
}

func (m *AllOfComponentMatcher) matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse {
	result := start
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if result == nil {
			result = &intsets.Sparse{}
		}
		if !set.entities.Has(entity) {
			result.Remove(entity)
			return result
		}
	}
	result.Insert(entity)
	return result
}

type AnyOfComponentMatcher struct {
	Components []reflect.Type
}

func (m *AnyOfComponentMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if result == nil {
			result = &intsets.Sparse{}
			result.Copy(set.entities)
		} else {
			result.UnionWith(set.entities)
		}
	}
	return result
}

func (m *AnyOfComponentMatcher) matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse {
	result := start
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if result == nil {
			result = &intsets.Sparse{}
		}
		if set.entities.Has(entity) {
			result.Insert(entity)
			return result
		}
	}
	result.Remove(entity)
	return result
}

type NoneOfComponentMatcher struct {
	Components []reflect.Type
}

func (m *NoneOfComponentMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
		result.Copy(storage.entities)
	}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result.DifferenceWith(set.entities)
	}
	return result
}

func (m *NoneOfComponentMatcher) matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
	}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if set.entities.Has(entity) {
			result.Remove(entity)
			return result
		}
	}
	result.Insert(entity)
	return result
}

type AllOfMatcher struct {
	Matchers []Matcher
}

func (m *AllOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	for _, mx := range m.Matchers {
		if result == nil {
			result = mx.match(nil, storage)
		} else {
			result.IntersectionWith(mx.match(result, storage))
		}
	}
	return result
}

func (m *AllOfMatcher) matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse {
	result := start
	for _, mx := range m.Matchers {
		if result == nil {
			result = mx.matchOne(nil, storage, entity)
		} else {
			result.IntersectionWith(mx.matchOne(result, storage, entity))
		}
	}
	return result
}

type AnyOfMatcher struct {
	Matchers []Matcher
}

func (m *AnyOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
	}
	for _, mx := range m.Matchers {
		result.UnionWith(mx.match(result, storage))
	}
	return result
}

func (m *AnyOfMatcher) matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
	}
	for _, mx := range m.Matchers {
		result.UnionWith(mx.matchOne(result, storage, entity))
	}
	return result
}

type NoneOfMatcher struct {
	Matchers []Matcher
}

func (m *NoneOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
		result.Copy(storage.entities)
	}
	for _, mx := range m.Matchers {
		result.DifferenceWith(mx.match(result, storage))
	}
	return result
}

func (m *NoneOfMatcher) matchOne(start *intsets.Sparse, storage *ComponentStorage, entity int) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
	}
	for _, mx := range m.Matchers {
		result.DifferenceWith(mx.matchOne(result, storage, entity))
	}
	return result
}

type Group struct {
	matcher       Matcher
	result        *intsets.Sparse
	entities      []int
	EntityAdded   chan int
	EntityRemoved chan int
}

func newGroup(matcher Matcher, storage *ComponentStorage) *Group {
	return &Group{
		matcher:       matcher,
		result:        matcher.match(nil, storage),
		EntityAdded:   make(chan int),
		EntityRemoved: make(chan int),
	}
}

func (g *Group) GetEntities() []int {
	// TODO: Caching.. But dirty later when running update
	if g.entities != nil {
		return g.entities
	}

	result := &intsets.Sparse{}
	result.Copy(g.result)

	entities := make([]int, result.Len())
	for i := 0; ; i++ {
		val := result.Min()
		if val == intsets.MaxInt {
			break
		}
		entities[i] = val
		result.Remove(val)
	}
	g.entities = entities
	return entities
}

func (g *Group) EvaluateEntity(entity int, storage *ComponentStorage) {
	go func() {
		if !storage.entities.Has(entity) {
			if g.result.Has(entity) {
				g.result.Remove(entity)
				// TODO : Better way to do this
				g.entities = nil
				g.EntityRemoved <- entity
			}
			return
		}
		result := &intsets.Sparse{}
		if g.result.Has(entity) {
			result.Insert(entity)
		}
		result = g.matcher.matchOne(result, storage, entity)

		if result.Has(entity) && g.result.Has(entity) {
			return
		}
		if !result.Has(entity) && !g.result.Has(entity) {
			return
		}

		if result.Has(entity) {
			g.result.Insert(entity)
			// TODO : Better way to do this
			g.entities = nil
			g.EntityAdded <- entity
		} else {
			g.result.Remove(entity)
			// TODO : Better way to do this
			g.entities = nil
			g.EntityRemoved <- entity
		}
	}()
}
