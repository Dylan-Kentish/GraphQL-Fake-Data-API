package api_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/api"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScehema(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User tests")
}

var _ = Describe("Users", func() {
	It("Invalid ID", func() {
		// Query
		query := `{user(id:"-1"){id,name,username}}`
		params := graphql.Params{Schema: api.UserSchema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		userJSON, _ := json.Marshal(result["user"])
		var user api.User
		json.Unmarshal(userJSON, &user)

		Expect(user).To(Equal(api.User{}))
	})

	userTests := make([]TableEntry, len(api.UserData))

	for i, user := range api.UserData {
		id, _ := strconv.Atoi(user.ID)
		userTests[i] = Entry(user.ID, id)
	}

	DescribeTable("Responds to query", func(id int) {
		// Query
		query := fmt.Sprintf(`{user(id:"%s"){id,name,username}}`, fmt.Sprint(id))
		params := graphql.Params{Schema: api.UserSchema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		userJSON, _ := json.Marshal(result["user"])
		var user api.User
		json.Unmarshal(userJSON, &user)

		Expect(user).To(Equal(api.UserData[id]))
	}, userTests)
})
