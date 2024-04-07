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
