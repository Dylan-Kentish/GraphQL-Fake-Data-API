package api_test

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/api"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

var _ = Describe("Photos", func() {
	data := api.NewData()
	testApi := api.NewAPI(data)

	photoTests := make([]TableEntry, len(data.Photos))
	for i, photo := range data.Photos {
		idString := fmt.Sprint(photo.ID)
		photoTests[i] = Entry(idString, photo.ID)
	}

	userTests := make([]TableEntry, len(data.Users))
	for i, user := range data.Users {
		idString := fmt.Sprint(user.ID)
		userTests[i] = Entry(idString, user.ID)
	}

	It("Invalid ID", func() {
		// Query
		query := `{photo(id:-1){id,albumid,description}}`
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var photo api.Photo
		convertTo(result["photo"], &photo)

		Expect(photo).To(Equal(api.Photo{}))
	})

	DescribeTable("Get photo by ID", func(id int) {
		// Query
		query := fmt.Sprintf(`{photo(id:%v){id,userid,description}}`, id)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var photo api.Photo
		convertTo(result["photo"], &photo)

		Expect(photo).To(Equal(data.Photos[id]))
	}, photoTests)

	DescribeTable("Get photo by albumID", func(albumId int) {
		// Query
		query := fmt.Sprintf(`{Photos(userid:%v){id,userid,description}}`, albumId)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var Photos []api.Photo
		convertTo(result["photos"], &Photos)

		expected := make([]api.Photo, 0)

		for _, photo := range data.Photos {
			if photo.AlbumID == albumId {
				expected = append(expected, photo)
			}
		}

		Expect(Photos).To(ContainElements(expected))
	}, userTests)

	It("Get all Photos", func() {
		// Query
		query := `{Photos{id,albumid,description}}`
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var Photos []api.Photo
		convertTo(result["Photos"], &Photos)

		Expect(Photos).To(ContainElements(maps.Values(data.Photos)))
	})

	It("Get limited Photos", func() {
		limit := 5
		// Query
		query := fmt.Sprintf(`{Photos(limit:%v){id,albumid,description}}`, limit)
		params := graphql.Params{Schema: testApi.Schema, RequestString: query}
		r := graphql.Do(params)
		Expect(r.Errors).To(BeEmpty())

		result := r.Data.(map[string]interface{})
		var Photos []api.Photo
		convertTo(result["Photos"], &Photos)

		Expect(Photos).To(HaveLen(limit))
	})

	Context("Bad Schema", func() {
		badQuery := graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"photo": &graphql.Field{
					Type:        testApi.PhotoType,
					Description: "photo by ID",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Description: "id of the photo",
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
			query := fmt.Sprintf(`{photo(id:0){%s}}`, field)
			params := graphql.Params{Schema: badSchema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(HaveLen(1))
			Expect(r.Errors[0].Message).To(Equal("source is not a api.Photo"))
		},
			Entry("id", "id"),
			Entry("albumid", "albumid"),
			Entry("description", "description"))
	})
})
