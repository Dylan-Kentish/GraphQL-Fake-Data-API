package tests

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/api"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

var _ = Describe("Albums", func() {
	data := api.NewData()
	testApi := api.NewAPI(data)

	albumTests := make([]TableEntry, len(data.Albums))
	for i, album := range data.Albums {
		idString := fmt.Sprint(album.ID)
		albumTests[i] = Entry(idString, album.ID)
	}

	userTests := make([]TableEntry, len(data.Users))
	for i, user := range data.Users {
		idString := fmt.Sprint(user.ID)
		userTests[i] = Entry(idString, user.ID)
	}

	It("Invalid ID", func() {
		// Query
		query := `{album(id:-1){id,userid,description}}`
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var album api.Album
		convertTo(result["album"], &album)

		Expect(album).To(Equal(api.Album{}))
	})

	DescribeTable("Get album by ID", func(id int) {
		// Query
		query := fmt.Sprintf(`{album(id:%v){id,userid,description}}`, id)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var album api.Album
		convertTo(result["album"], &album)

		Expect(album).To(Equal(data.Albums[id]))
	}, albumTests)

	DescribeTable("Get album by userID", func(userId int) {
		// Query
		query := fmt.Sprintf(`{albums(userid:%v){id,userid,description}}`, userId)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var albums []api.Album
		convertTo(result["albums"], &albums)

		expected := make([]api.Album, 0)

		for _, album := range data.Albums {
			if album.UserID == userId {
				expected = append(expected, album)
			}
		}

		Expect(albums).To(ContainElements(expected))
	}, userTests)

	It("Get all albums", func() {
		// Query
		query := `{albums{id,userid,description}}`
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var albums []api.Album
		convertTo(result["albums"], &albums)

		Expect(albums).To(ContainElements(maps.Values(data.Albums)))
	})

	DescribeTable("Get album photos", func(id int) {
		query := fmt.Sprintf(`{album(id:%v){id,photos{id,albumid,description}}}`, id)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var album api.Album
		convertTo(result["album"], &album)

		expected := make([]api.Photo, 0)

		for _, photo := range data.Photos {
			if photo.AlbumID == id {
				expected = append(expected, photo)
			}
		}

		Expect(album.Photos).To(ContainElements(expected))
	}, userTests)

	It("Get limited albums", func() {
		limit := 5
		// Query
		query := fmt.Sprintf(`{albums(limit:%v){id,userid,description}}`, limit)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var albums []api.Album
		convertTo(result["albums"], &albums)

		Expect(albums).To(HaveLen(limit))
	})

	Context("Bad Schema", func() {
		badQuery := graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"album": &graphql.Field{
					Type:        testApi.AlbumType,
					Description: "Album by ID",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Description: "id of the album",
							Type:        graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// Wrong type
						return api.User{}, nil
					},
				},
			},
		})

		badSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
			Query: badQuery,
		})

		DescribeTable("Reterns err when resolving fields", func(field string) {
			// Query
			query := fmt.Sprintf(`{album(id:0){%s}}`, field)
			params := graphql.Params{Schema: badSchema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(HaveLen(1))
			Expect(r.Errors[0].Message).To(Equal("source is not a api.Album"))
		},
			Entry("id", "id"),
			Entry("userid", "userid"),
			Entry("description", "description"),
			Entry("photos", "photos{id}"))
	})
})
