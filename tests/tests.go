package tests

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/onsi/ginkgo/v2"
	"golang.org/x/exp/maps"
)

func ConvertFieldDefinitionToQueryString(item *graphql.FieldDefinition) ginkgo.TableEntry {
	value := item.Name

	if list, ok := item.Type.(*graphql.List); ok {
		if obj, ok := list.OfType.(*graphql.Object); ok {
			subFieldsMap := obj.Fields()
			if len(subFieldsMap) > 0 {
				subFieldsValues := maps.Values(subFieldsMap)
				value += fmt.Sprintf("{%s}", subFieldsValues[0].Name)
			}
		}
	}

	return ginkgo.Entry(item.Name, value)
}
