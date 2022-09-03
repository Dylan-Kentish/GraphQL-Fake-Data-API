package utils

import "encoding/json"

// Transform applies the transformation to every item in the slice m and returns the result.
func Transform[M []T, T any, K any](m M, transformation func(item T) K) []K {
	r := make([]K, 0, len(m))
	for _, v := range m {
		r = append(r, transformation(v))
	}
	return r
}

// Converts the provided interface into struct T via Json.
func ConvertTo[T any](in interface{}, out *T) {
	bytes, _ := json.Marshal(in)
	json.Unmarshal(bytes, out)
}
