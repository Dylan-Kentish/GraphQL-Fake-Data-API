package testUtils

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/onsi/ginkgo/v2"
	"golang.org/x/exp/maps"
)

func ConvertFieldDefinitionToTableEntry(item *graphql.FieldDefinition) ginkgo.TableEntry {
	value := convertFieldDefinitionToQueryString(item)

	return ginkgo.Entry(item.Name, value)
}

func convertFieldDefinitionToQueryString(item *graphql.FieldDefinition) string {
	value := item.Name

	if list, ok := item.Type.(*graphql.List); ok {
		if obj, ok := list.OfType.(*graphql.Object); ok {
			subFieldsMap := obj.Fields()
			if len(subFieldsMap) > 0 {
				field := maps.Values(subFieldsMap)[0]
				value += fmt.Sprintf("{%s}", convertFieldDefinitionToQueryString(field))
			}
		}
	}

	return value
}

func GetData[T any](r *graphql.Result, key string) T {
	result := r.Data.(map[string]interface{})
	return convertTo[T](result[key])
}

// Converts the provided interface into struct T via Json.
func convertTo[T any](in interface{}) T {
	bytes, _ := json.Marshal(in)
	var out T
	json.Unmarshal(bytes, &out)
	return out
}
