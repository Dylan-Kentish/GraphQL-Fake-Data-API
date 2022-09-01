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
	"golang.org/x/exp/maps"
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
		var user api.User
		convertTo(result["user"], &user)

		Expect(user).To(Equal(api.User{}))
	})

	userTests := make([]TableEntry, len(api.UserData))

	for i, user := range api.UserData {
		id, _ := strconv.Atoi(user.ID)
		userTests[i] = Entry(user.ID, id)
	}

	DescribeTable("Get user by ID", func(id int) {
		// Query
		query := fmt.Sprintf(`{user(id:"%s"){id,name,username}}`, fmt.Sprint(id))
		params := graphql.Params{Schema: api.UserSchema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		Expect(user).To(Equal(api.UserData[id]))
	}, userTests)

	It("Get all users", func() {
		// Query
		query := `{users{id,name,username}}`
		params := graphql.Params{Schema: api.UserSchema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var users []api.User
		convertTo(result["users"], &users)

		Expect(users).To(ContainElements(maps.Values(api.UserData)))
	})
})

func convertTo[T any](in interface{}, out *T) {
	bytes, _ := json.Marshal(in)
	json.Unmarshal(bytes, out)
}
