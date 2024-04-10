package engine

import (
	"golang.org/x/tools/container/intsets"
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
		expected []int
	}{
		{
			name:     "AllOfComponentMatcher_1",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{}},
			expected: []int{},
		},
		{
			name:     "AllOfComponentMatcher_2",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			expected: []int{0, 2},
		},
		{
			name:     "AllOfComponentMatcher_3",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			expected: []int{1, 2},
		},
		{
			name:     "AllOfComponentMatcher_4",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeA, typeB}},
			expected: []int{2},
		},
		{
			name:     "AllOfComponentMatcher_5",
			matcher:  &AllOfComponentMatcher{Components: []reflect.Type{typeA, typeB, typeC}},
			expected: []int{},
		},
		{
			name:     "AnyOfComponentMatcher_1",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{}},
			expected: []int{},
		},
		{
			name:     "AnyOfComponentMatcher_2",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{typeA}},
			expected: []int{0, 2},
		},
		{
			name:     "AnyOfComponentMatcher_3",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{typeA, typeB}},
			expected: []int{0, 1, 2},
		},
		{
			name:     "AnyOfComponentMatcher_4",
			matcher:  &AnyOfComponentMatcher{Components: []reflect.Type{typeA, typeB, typeC}},
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "NoneOfComponentMatcher_1",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{}},
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "NoneOfComponentMatcher_2",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeA}},
			expected: []int{1, 3},
		},
		{
			name:     "NoneOfComponentMatcher_3",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			expected: []int{0, 3},
		},
		{
			name:     "NoneOfComponentMatcher_4",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeC}},
			expected: []int{0, 1, 2},
		},
		{
			name:     "NoneOfComponentMatcher_5",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeA, typeB}},
			expected: []int{3},
		},
		{
			name:     "NoneOfComponentMatcher_6",
			matcher:  &NoneOfComponentMatcher{Components: []reflect.Type{typeA, typeB, typeC}},
			expected: []int{},
		},
		{
			name: "AllOfMatcher_1",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{}},
			}},
			expected: []int{},
		},
		{
			name: "AllOfMatcher_2",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			}},
			expected: []int{0, 2},
		},
		{
			name: "AllOfMatcher_3",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{2},
		},
		{
			name: "AllOfMatcher_4",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{0},
		},
		{
			name: "AllOfMatcher_5",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeC}},
			}},
			expected: []int{},
		},
		{
			name: "AllOfMatcher_6",
			matcher: &AllOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{2},
		},
		{
			name: "AnyOfMatcher_1",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{}},
			}},
			expected: []int{},
		},
		{
			name: "AnyOfMatcher_2",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			}},
			expected: []int{0, 2},
		},
		{
			name: "AnyOfMatcher_3",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{0, 1, 2},
		},
		{
			name: "AnyOfMatcher_4",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{0, 2, 3},
		},
		{
			name: "AnyOfMatcher_5",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeC}},
			}},
			expected: []int{0, 2, 3},
		},
		{
			name: "AnyOfMatcher_6",
			matcher: &AnyOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{0, 1, 2},
		},
		{
			name: "NoneOfMatcher_1",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{}},
			}},
			expected: []int{0, 1, 2, 3},
		},
		{
			name: "NoneOfMatcher_2",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
			}},
			expected: []int{1, 3},
		},
		{
			name: "NoneOfMatcher_3",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AllOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{3},
		},
		{
			name: "NoneOfMatcher_4",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&NoneOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{1},
		},
		{
			name: "NoneOfMatcher_5",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeC}},
			}},
			expected: []int{1},
		},
		{
			name: "NoneOfMatcher_6",
			matcher: &NoneOfMatcher{Matchers: []Matcher{
				&AllOfComponentMatcher{Components: []reflect.Type{typeA}},
				&AnyOfComponentMatcher{Components: []reflect.Type{typeB}},
			}},
			expected: []int{3},
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
				assert.ElementsMatch(t, []int{expected}, actual)
			}
		})
	}
}

// Helper function to convert intsets.Sparse to slice
func toSlice(s *intsets.Sparse) []int {
	result := &intsets.Sparse{}
	result.Copy(s)

	entities := make([]int, result.Len())
	for i := 0; ; i++ {
		val := result.Min()
		if val == intsets.MaxInt {
			break
		}
		entities[i] = val
		result.Remove(val)
	}
	return entities
}

// Mock component types for testing
type ComponentA struct{}
type ComponentB struct{}
type ComponentC struct{}
