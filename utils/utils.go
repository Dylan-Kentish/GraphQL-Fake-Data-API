package utils

func TryLimitIfPresent[S []T, T any](s S, Args map[string]interface{}) []T {
	if limit, exists := Args["limit"].(int); exists &&
		limit < len(s) {
		return s[:limit]
	}

	return s
}
