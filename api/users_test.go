package api_test

import (
	"encoding/json"
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/api"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

var _ = Describe("Users", func() {
	var testApi *api.API
	var userTests []TableEntry

	BeforeEach(func() {
		testApi = api.NewAPI()

		userTests = make([]TableEntry, len(testApi.Data.Users))
		for i, user := range testApi.Data.Users {
			idString := fmt.Sprint(user.ID)
			userTests[i] = Entry(idString, user.ID)
		}
	})

	It("Invalid ID", func() {
		// Query
		query := `{user(id:-1){id,name,username}}`
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		Expect(user).To(Equal(api.User{}))
	})

	DescribeTable("Get user by ID", func(id int) {
		// Query
		query := fmt.Sprintf(`{user(id:%s){id,name,username}}`, fmt.Sprint(id))
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		Expect(user).To(Equal(testApi.Data.Users[id]))
	}, userTests)

	It("Get all users", func() {
		// Query
		query := `{users{id,name,username}}`
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var users []api.User
		convertTo(result["users"], &users)

		Expect(users).To(ContainElements(maps.Values(testApi.Data.Users)))
	})

	DescribeTable("Get user albums", func(id int) {
		query := fmt.Sprintf(`{user(id:%v){id,albums{id,userid,description}}}`, id)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		expected := make([]api.Album, 0)

		for _, album := range testApi.Data.Albums {
			if album.UserID == id {
				expected = append(expected, album)
			}
		}

		Expect(user.Albums).To(ContainElements(expected))
	}, userTests)

	It("Get limited users", func() {
		limit := 5
		// Query
		query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var users []api.User
		convertTo(result["users"], &users)

		Expect(users).To(HaveLen(limit))
	})
})

func convertTo[T any](in interface{}, out *T) {
	bytes, _ := json.Marshal(in)
	json.Unmarshal(bytes, out)
}
