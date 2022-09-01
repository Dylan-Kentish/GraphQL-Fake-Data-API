package internal_test

import (
	"encoding/json"
	"testing"

	. "github.com/Dylan-Kentish/GraphQLFakeDataAPI/internal"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScehema(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scehema tests")
}

var _ = Describe("Hello World", Label("Main"), func() {
	It("Responds to query", func() {
		// Query
		query := `
		{
			hello
		}
		`
		params := graphql.Params{Schema: Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		bJSON, _ := json.Marshal(r.Data)
		var result map[string]interface{}
		json.Unmarshal(bJSON, &result)

		Expect(result["hello"]).To(Equal("world"))
	})
})
