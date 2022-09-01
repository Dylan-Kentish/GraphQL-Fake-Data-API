package internal_test

import (
	"testing"

	. "github.com/Dylan-Kentish/GraphQLFakeDataAPI/internal"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSchema(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Schema tests")
}

var _ = Describe("Hello World", func() {
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

		result := r.Data.(map[string]interface{})

		Expect(result["hello"]).To(Equal("world"))
	})
})
