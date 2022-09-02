package api_test

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/api"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

var _ = Describe("Albums", func() {
	It("Invalid ID", func() {
		// Query
		query := `{album(id:-1){id,userid,description}}`
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var album api.Album
		convertTo(result["album"], &album)

		Expect(album).To(Equal(api.Album{}))
	})

	albumTests := make([]TableEntry, len(api.Data.Albums))

	for i, album := range api.Data.Albums {
		idString := fmt.Sprint(album.ID)
		albumTests[i] = Entry(idString, album.ID)
	}

	DescribeTable("Get album by ID", func(id int) {
		// Query
		query := fmt.Sprintf(`{album(id:%s){id,userid,description}}`, fmt.Sprint(id))
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var album api.Album
		convertTo(result["album"], &album)

		Expect(album).To(Equal(api.Data.Albums[id]))
	}, albumTests)

	userTests := make([]TableEntry, len(api.Data.Users))

	for i, user := range api.Data.Users {
		idString := fmt.Sprint(user.ID)
		userTests[i] = Entry(idString, user.ID)
	}

	DescribeTable("Get album by userID", func(userId int) {
		// Query
		query := fmt.Sprintf(`{albums(userid:%v){id,userid,description}}`, userId)
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var albums []api.Album
		convertTo(result["albums"], &albums)

		expected := make([]api.Album, 0)

		for _, album := range api.Data.Albums {
			if album.UserID == userId {
				expected = append(expected, album)
			}
		}

		Expect(albums).To(ContainElements(expected))
	}, userTests)

	It("Get all albums", func() {
		// Query
		query := `{albums{id,userid,description}}`
		params := graphql.Params{Schema: api.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var albums []api.Album
		convertTo(result["albums"], &albums)

		Expect(albums).To(ContainElements(maps.Values(api.Data.Albums)))
	})
})
