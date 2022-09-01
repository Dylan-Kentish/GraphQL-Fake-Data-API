package utils

import (
	"encoding/json"
)

func As[T any](obj interface{}, result *T) {
	bytes, _ := json.Marshal(obj)
	json.Unmarshal(bytes, result)
}
