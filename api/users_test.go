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
	It("Invalid ID", func() {
		// Query
		query := `{user(id:-1){id,name,username}}`
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		Expect(user).To(Equal(api.User{}))
	})

	userTests := make([]TableEntry, len(api.Data.Users))

	for i, user := range api.Data.Users {
		idString := fmt.Sprint(user.ID)
		userTests[i] = Entry(idString, user.ID)
	}

	DescribeTable("Get user by ID", func(id int) {
		// Query
		query := fmt.Sprintf(`{user(id:%s){id,name,username}}`, fmt.Sprint(id))
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		Expect(user).To(Equal(api.Data.Users[id]))
	}, userTests)

	It("Get all users", func() {
		// Query
		query := `{users{id,name,username}}`
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var users []api.User
		convertTo(result["users"], &users)

		Expect(users).To(ContainElements(maps.Values(api.Data.Users)))
	})

	DescribeTable("Get user albums", func(id int) {
		query := fmt.Sprintf(`{user(userid:%s){id,albums}}`, fmt.Sprint(id))
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var user api.User
		convertTo(result["user"], &user)

		expected := make([]api.Album, 0)

		for _, album := range api.Data.Albums {
			if album.UserID == id {
				expected = append(expected, album)
			}
		}

		Expect(user.Albums).To(ContainElements(expected))
	})

	It("Get limited users", func() {
		limit := 5
		// Query
		query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
		params := graphql.Params{Schema: api.Schema, RequestString: query}
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
