package utils

import (
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// Transform applies the transformation to every item in the slice m and returns the result.
func Transform[S []T, T any, U any](m S, transformation func(item T) U) []U {
	r := make([]U, 0, len(m))
	for _, v := range m {
		r = append(r, transformation(v))
	}
	return r
}

// TransformValues applies the transformation to every value in the map m and returns the resulting slice.
func TransformValues[M ~map[T]U, T constraints.Ordered, U any, V any](m M, transformation func(item U) V) []V {
	return Transform(OrderedValues(m), transformation)
}

// Returns all items that match the condition.
func Where[S []T, T any](m S, condition func(item T) bool) S {
	r := make(S, 0)
	for _, item := range m {
		if condition(item) {
			r = append(r, item)
		}
	}
	return r
}

// Returns all values that match the condition.
func ValuesWhere[M map[T]U, T constraints.Ordered, U any](m M, condition func(item U) bool) []U {
	return Where(OrderedValues(m), condition)
}

func TryLimitIfPresent[S []T, T any](s S, Args map[string]interface{}) []T {
	if limit, exists := Args["limit"].(int); exists &&
		limit < len(s) {
		return s[:limit]
	}

	return s
}

// OrderedValues returns the values sorted by Key
func OrderedValues[M ~map[T]U, T constraints.Ordered, U any](m M) []U {
	keys := maps.Keys(m)

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	result := make([]U, 0)

	for _, key := range keys {
		result = append(result, m[key])
	}

	return result
}
