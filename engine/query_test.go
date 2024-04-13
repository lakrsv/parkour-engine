package engine

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test matching functionality
func TestMatcherMatch(t *testing.T) {
	// Mock component storage
	storage := NewComponentStorage()

	// Define some component types
	typeA := reflect.TypeOf(ComponentA{})
	typeB := reflect.TypeOf(ComponentB{})
	typeC := reflect.TypeOf(ComponentC{})

	// Add some entities to storage
	storage.createEntity(ComponentA{})
	storage.createEntity(ComponentB{})
	storage.createEntity(ComponentA{}, ComponentB{})
	storage.createEntity(ComponentC{})

	// Define test cases
	testCases := []struct {
		name     string
		matcher  Matcher
		expected []uint32
	}{
		{
			name:     "AllOfComponentMatcher_1",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{}},
			expected: []uint32{},
		},
		{
			name:     "AllOfComponentMatcher_2",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			expected: []uint32{0, 2},
		},
		{
			name:     "AllOfComponentMatcher_3",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			expected: []uint32{1, 2},
		},
		{
			name:     "AllOfComponentMatcher_4",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeA, typeB}},
			expected: []uint32{2},
		},
		{
			name:     "AllOfComponentMatcher_5",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeA, typeB, typeC}},
			expected: []uint32{},
		},
		{
			name:     "AnyOfComponentMatcher_1",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{}},
			expected: []uint32{},
		},
		{
			name:     "AnyOfComponentMatcher_2",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{typeA}},
			expected: []uint32{0, 2},
		},
		{
			name:     "AnyOfComponentMatcher_3",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{typeA, typeB}},
			expected: []uint32{0, 1, 2},
		},
		{
			name:     "AnyOfComponentMatcher_4",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{typeA, typeB, typeC}},
			expected: []uint32{0, 1, 2, 3},
		},
		{
			name:     "NoneOfComponentMatcher_1",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{}},
			expected: []uint32{0, 1, 2, 3},
		},
		{
			name:     "NoneOfComponentMatcher_2",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeA}},
			expected: []uint32{1, 3},
		},
		{
			name:     "NoneOfComponentMatcher_3",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			expected: []uint32{0, 3},
		},
		{
			name:     "NoneOfComponentMatcher_4",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeC}},
			expected: []uint32{0, 1, 2},
		},
		{
			name:     "NoneOfComponentMatcher_5",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeA, typeB}},
			expected: []uint32{3},
		},
		{
			name:     "NoneOfComponentMatcher_6",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeA, typeB, typeC}},
			expected: []uint32{},
		},
		{
			name: "AllOfMatcher_1",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{}},
			}},
			expected: []uint32{},
		},
		{
			name: "AllOfMatcher_2",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			}},
			expected: []uint32{2, 0},
		},
		{
			name: "AllOfMatcher_3",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{2},
		},
		{
			name: "AllOfMatcher_4",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{0},
		},
		{
			name: "AllOfMatcher_5",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeC}},
			}},
			expected: []uint32{},
		},
		{
			name: "AllOfMatcher_6",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{2},
		},
		{
			name: "AnyOfMatcher_1",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{}},
			}},
			expected: []uint32{},
		},
		{
			name: "AnyOfMatcher_2",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			}},
			expected: []uint32{0, 2},
		},
		{
			name: "AnyOfMatcher_3",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{0, 1, 2},
		},
		{
			name: "AnyOfMatcher_4",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{0, 2, 3},
		},
		{
			name: "AnyOfMatcher_5",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeC}},
			}},
			expected: []uint32{0, 2, 3},
		},
		{
			name: "AnyOfMatcher_6",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{0, 1, 2},
		},
		{
			name: "NoneOfMatcher_1",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{}},
			}},
			expected: []uint32{0, 1, 2, 3},
		},
		{
			name: "NoneOfMatcher_2",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			}},
			expected: []uint32{1, 3},
		},
		{
			name: "NoneOfMatcher_3",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{3},
		},
		{
			name: "NoneOfMatcher_4",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{1},
		},
		{
			name: "NoneOfMatcher_5",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeC}},
			}},
			expected: []uint32{1},
		},
		{
			name: "NoneOfMatcher_6",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []uint32{3},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test match function
			result := tc.matcher.match(storage)
			actual := toSlice(result)
			assert.ElementsMatch(t, tc.expected, actual)
			// Test match one function
			for _, expected := range tc.expected {
				result := tc.matcher.matchOne(storage, expected)
				actual := toSlice(result)
				assert.ElementsMatch(t, []uint32{expected}, actual)
			}
		})
	}
}

// Helper function to convert intsets.Sparse to slice
func toSlice(s *SparseSet[uint32]) []uint32 {
	entities := make([]uint32, s.Len())
	iterator := s.Iterator()
	idx := 0
	for {
		id, _, ok := iterator.Next()
		if !ok {
			break
		}
		entities[idx] = id
		idx += 1
	}
	return entities
}

// Mock component types for testing
type ComponentA struct{}
type ComponentB struct{}
type ComponentC struct{}
