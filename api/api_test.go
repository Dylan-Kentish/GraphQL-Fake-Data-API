package api

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/data"
	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/testUtils"
	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/utils"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

var _ = Describe("Api", func() {
	testData := testUtils.NewTestData()
	testApi := NewAPI(testData)

	photoTests := make([]TableEntry, len(testData.Photos))
	for i, photo := range testData.Photos {
		idString := fmt.Sprint(photo.ID)
		photoTests[i] = Entry(idString, photo.ID)
	}

	albumTests := make([]TableEntry, len(testData.Albums))
	for i, album := range testData.Albums {
		idString := fmt.Sprint(album.ID)
		albumTests[i] = Entry(idString, album.ID)
	}

	userTests := make([]TableEntry, len(testData.Users))
	for i, user := range testData.Users {
		idString := fmt.Sprint(user.ID)
		userTests[i] = Entry(idString, user.ID)
	}

	Context("Albums", func() {
		It("Invalid ID", func() {
			// Query
			query := `{album(id:-1){id,userid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var album data.Album
			utils.ConvertTo(result["album"], &album)

			Expect(album).To(Equal(data.Album{}))
		})

		DescribeTable("Get album by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{album(id:%v){id,userid,description}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var album data.Album
			utils.ConvertTo(result["album"], &album)

			Expect(album).To(Equal(testData.Albums[id]))
		}, albumTests)

		DescribeTable("Get album by userID", func(userId int) {
			// Query
			query := fmt.Sprintf(`{albums(userid:%v){id,userid,description}}`, userId)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var albums []data.Album
			utils.ConvertTo(result["albums"], &albums)

			expected := make([]data.Album, 0)

			for _, album := range testData.Albums {
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
			var albums []data.Album
			utils.ConvertTo(result["albums"], &albums)

			Expect(albums).To(ContainElements(maps.Values(testData.Albums)))
		})

		DescribeTable("Get album photos", func(id int) {
			query := fmt.Sprintf(`{album(id:%v){id,photos{id,albumid,description}}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var album data.Album
			utils.ConvertTo(result["album"], &album)

			expected := make([]data.Photo, 0)

			for _, photo := range testData.Photos {
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
			var albums []data.Album
			utils.ConvertTo(result["albums"], &albums)

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
							return data.User{}, nil
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

	Context("Photos", func() {
		It("Invalid ID", func() {
			// Query
			query := `{photo(id:-1){id,albumid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var photo data.Photo
			utils.ConvertTo(result["photo"], &photo)

			Expect(photo).To(Equal(data.Photo{}))
		})

		DescribeTable("Get photo by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{photo(id:%v){id,albumid,description}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var photo data.Photo
			utils.ConvertTo(result["photo"], &photo)

			Expect(photo).To(Equal(testData.Photos[id]))
		}, photoTests)

		DescribeTable("Get photo by albumID", func(albumId int) {
			// Query
			query := fmt.Sprintf(`{photos(albumid:%v){id,albumid,description}}`, albumId)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var Photos []data.Photo
			utils.ConvertTo(result["photos"], &Photos)

			expected := make([]data.Photo, 0)

			for _, photo := range testData.Photos {
				if photo.AlbumID == albumId {
					expected = append(expected, photo)
				}
			}

			Expect(Photos).To(ContainElements(expected))
		}, userTests)

		It("Get all Photos", func() {
			// Query
			query := `{photos{id,albumid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var Photos []data.Photo
			utils.ConvertTo(result["photos"], &Photos)

			Expect(Photos).To(ContainElements(maps.Values(testData.Photos)))
		})

		It("Get limited Photos", func() {
			limit := 5
			// Query
			query := fmt.Sprintf(`{photos(limit:%v){id,albumid,description}}`, limit)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var Photos []data.Photo
			utils.ConvertTo(result["photos"], &Photos)

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
							return data.User{}, nil
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

	Context("Users", func() {
		It("Invalid ID", func() {
			// Query
			query := `{user(id:-1){id,name,username}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var user data.User
			utils.ConvertTo(result["user"], &user)

			Expect(user).To(Equal(data.User{}))
		})

		DescribeTable("Get user by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{user(id:%v){id,name,username}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var user data.User
			utils.ConvertTo(result["user"], &user)

			Expect(user).To(Equal(testData.Users[id]))
		}, userTests)

		It("Get all users", func() {
			// Query
			query := `{users{id,name,username}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var users []data.User
			utils.ConvertTo(result["users"], &users)

			Expect(users).To(ContainElements(maps.Values(testData.Users)))
		})

		DescribeTable("Get user albums", func(id int) {
			query := fmt.Sprintf(`{user(id:%v){id,albums{id,userid,description}}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			result := r.Data.(map[string]interface{})
			var user data.User
			utils.ConvertTo(result["user"], &user)

			expected := make([]data.Album, 0)

			for _, album := range testData.Albums {
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
			var users []data.User
			utils.ConvertTo(result["users"], &users)

			Expect(users).To(HaveLen(limit))
		})

		Context("Bad Schema", func() {
			badQuery := graphql.NewObject(graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					"user": &graphql.Field{
						Type:        testApi.UserType,
						Description: "User by ID",
						Args: graphql.FieldConfigArgument{
							"id": &graphql.ArgumentConfig{
								Description: "id of the user",
								Type:        graphql.NewNonNull(graphql.Int),
							},
						},
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							// Wrong type
							return data.Album{}, nil
						},
					},
					"users": &graphql.Field{
						Type:        graphql.NewList(testApi.UserType),
						Description: "All users",
						Args: graphql.FieldConfigArgument{
							"limit": &graphql.ArgumentConfig{
								Description: "limit the number of users",
								Type:        graphql.Int,
							},
						},
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							// Wrong type
							albums := make([]data.Album, 0)
							return albums, nil
						},
					},
				},
			})

			badSchema, _ := graphql.NewSchema(graphql.SchemaConfig{
				Query: badQuery,
			})

			DescribeTable("Reterns err when resolving fields", func(field string) {
				// Query
				query := fmt.Sprintf(`{user(id:0){%s}}`, field)
				params := graphql.Params{Schema: badSchema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not a api.User"))
			},
				Entry("id", "id"),
				Entry("name", "name"),
				Entry("username", "username"),
				Entry("albums", "albums{id}"))
		})
	})
})
