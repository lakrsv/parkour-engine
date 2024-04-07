package engine

import (
	"golang.org/x/tools/container/intsets"
	"reflect"
)

type Matcher interface {
	match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse
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
			result.Copy(&set.entities)
		} else {
			result.IntersectionWith(&set.entities)
		}
	}
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
			result.Copy(&set.entities)
		} else {
			result.UnionWith(&set.entities)
		}
	}
	return result
}

type NoneOfComponentMatcher struct {
	Components []reflect.Type
}

func (m *NoneOfComponentMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
		result.Copy(&storage.entities)
	}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result.DifferenceWith(&set.entities)
	}
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

type NoneOfMatcher struct {
	Matchers []Matcher
}

func (m *NoneOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	if result == nil {
		result = &intsets.Sparse{}
		result.Copy(&storage.entities)
	}
	for _, mx := range m.Matchers {
		result.DifferenceWith(mx.match(result, storage))
	}
	return result
}

type Group struct {
	matcher  Matcher
	result   *intsets.Sparse
	entities []int
}

func newGroup(matcher Matcher, storage *ComponentStorage) *Group {
	return &Group{
		matcher: matcher,
		result:  matcher.match(nil, storage),
	}
}

func (g *Group) GetEntities() []int {
	// TODO: Caching.. But dirty later when running update
	if g.entities != nil {
		return g.entities
	}

	result := intsets.Sparse{}
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
