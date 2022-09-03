package utils

import "golang.org/x/exp/maps"

// Transform applies the transformation to every item in the slice m and returns the result.
func Transform[M []T, T any, K any](m M, transformation func(item T) K) []K {
	r := make([]K, 0, len(m))
	for _, v := range m {
		r = append(r, transformation(v))
	}
	return r
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
