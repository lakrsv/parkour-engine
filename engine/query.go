package engine

import (
	"golang.org/x/tools/container/intsets"
	"reflect"
)

type Matcher interface {
	match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse
}

type AllOfMatcher struct {
	Components []reflect.Type
}

func (m *AllOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if result.IsEmpty() {
			result.Copy(&set.entities)
		} else {
			result.IntersectionWith(&set.entities)
		}
	}
	return result
}

type AnyOfMatcher struct {
	Components []reflect.Type
}

func (m *AnyOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		if result.IsEmpty() {
			result.Copy(&set.entities)
		} else {
			result.UnionWith(&set.entities)
		}
	}
	return result
}

type NoneOfMatcher struct {
	Components []reflect.Type
}

func (m *NoneOfMatcher) match(start *intsets.Sparse, storage *ComponentStorage) *intsets.Sparse {
	result := start
	if result.IsEmpty() {
		result.Copy(&storage.entities)
	}
	for _, t := range m.Components {
		set := storage.getComponentSet(t)
		result.DifferenceWith(&set.entities)
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
		result:  matcher.match(&intsets.Sparse{}, storage),
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
