package utils

import (
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
func TransformValues[M ~map[T]U, T comparable, U any, V any](m M, transformation func(item U) V) []V {
	return Transform(maps.Values(m), transformation)
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
func ValuesWhere[M map[T]U, T comparable, U any](m M, condition func(item U) bool) []U {
	return Where(maps.Values(m), condition)
}

func TryLimitIfPresent[S []T, T any](s S, Args map[string]interface{}) []T {
	if limit, exists := Args["limit"].(int); exists &&
		limit < len(s) {
		return s[:limit]
	}

	return s
}
