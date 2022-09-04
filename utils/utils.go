package utils

import (
	"golang.org/x/exp/maps"
)

// Transform applies the transformation to every item in the slice m and returns the result.
func Transform[T []U, U any, V any](m T, transformation func(item U) V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, transformation(v))
	}
	return r
}

// TransformValues applies the transformation to every value in the map m and returns the resulting slice.
func TransformValues[T ~map[W]U, U any, V any, W comparable](m T, transformation func(item U) V) []V {
	return Transform(maps.Values(m), transformation)
}

// Returns all items that match the condition.
func Where[M []T, T any](m M, condition func(item T) bool) M {
	r := make(M, 0)
	for _, item := range m {
		if condition(item) {
			r = append(r, item)
		}
	}
	return r
}

// Returns all values that match the condition.
func ValuesWhere[M map[K]T, T any, K comparable](m M, condition func(item T) bool) []T {
	return Where(maps.Values(m), condition)
}

func TryLimitIfPresent[S []T, T any](s S, Args map[string]interface{}) []T {
	if limit, exists := Args["limit"].(int); exists &&
		limit < len(s) {
		return s[:limit]
	}

	return s
}
