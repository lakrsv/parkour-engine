package engine

import (
	"golang.org/x/tools/container/intsets"
	"reflect"
)

type Matcher interface {
	match(storage *ComponentStorage) *intsets.Sparse
	matchOne(storage *ComponentStorage, entity int) *intsets.Sparse
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

func (m *AllOfComponentMatcher) match(storage *ComponentStorage) *intsets.Sparse {
	result := &intsets.Sparse{}
	for i, t := range m.Components {
		set := storage.getComponentSet(t)
		if result.IsEmpty() && i == 0 {
			result.Copy(set.entities)
		} else {
			result.IntersectionWith(set.entities)
		}
	}
	return result
}

func (m *AnyOfComponentMatcher) match(storage *ComponentStorage) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result.UnionWith(set.entities)
	}
	return result
}

func (m *NoneOfComponentMatcher) match(storage *ComponentStorage) *intsets.Sparse {
	result := &intsets.Sparse{}
	result.Copy(storage.entities)
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result.DifferenceWith(set.entities)
	}
	return result
}

func (m *AllOfMatcher) match(storage *ComponentStorage) *intsets.Sparse {
	result := &intsets.Sparse{}
	for i, mx := range m.Matchers {
		if i == 0 {
			result = mx.match(storage)
		} else {
			result.IntersectionWith(mx.match(storage))
		}
	}
	return result
}

func (m *AnyOfMatcher) match(storage *ComponentStorage) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, mx := range m.Matchers {
		result.UnionWith(mx.match(storage))
	}
	return result
}

func (m *NoneOfMatcher) match(storage *ComponentStorage) *intsets.Sparse {
	result := &intsets.Sparse{}
	result.Copy(storage.entities)
	for _, mx := range m.Matchers {
		result.DifferenceWith(mx.match(storage))
	}
	return result
}

func (m *AllOfComponentMatcher) matchOne(storage *ComponentStorage, entity int) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if !set.entities.Has(entity) {
			result.Remove(entity)
			return result
		}
	}
	result.Insert(entity)
	return result
}

func (m *AnyOfComponentMatcher) matchOne(storage *ComponentStorage, entity int) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if set.entities.Has(entity) {
			result.Insert(entity)
			return result
		}
	}
	result.Remove(entity)
	return result
}

func (m *NoneOfComponentMatcher) matchOne(storage *ComponentStorage, entity int) *intsets.Sparse {
	result := &intsets.Sparse{}
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

func (m *AllOfMatcher) matchOne(storage *ComponentStorage, entity int) *intsets.Sparse {
	result := &intsets.Sparse{}
	for i, mx := range m.Matchers {
		if i == 0 {
			result = mx.matchOne(storage, entity)
		} else {
			result.IntersectionWith(mx.matchOne(storage, entity))
		}
	}
	return result
}

func (m *AnyOfMatcher) matchOne(storage *ComponentStorage, entity int) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, mx := range m.Matchers {
		result.UnionWith(mx.matchOne(storage, entity))
	}
	return result
}

func (m *NoneOfMatcher) matchOne(storage *ComponentStorage, entity int) *intsets.Sparse {
	result := &intsets.Sparse{}
	for _, mx := range m.Matchers {
		result.DifferenceWith(mx.matchOne(storage, entity))
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
		result:        matcher.match(storage),
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
		result := g.matcher.matchOne(storage, entity)

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
