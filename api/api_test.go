package api

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/data"
	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/utils"
	"github.com/graphql-go/graphql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	testData := NewTestData()

	photoTests := utils.TransformValues(testData.Photos, func(v data.Photo) TableEntry {
		return Entry(fmt.Sprint(v.ID), v.ID)
	})

	albumTests := utils.TransformValues(testData.Albums, func(v data.Album) TableEntry {
		return Entry(fmt.Sprint(v.ID), v.ID)
	})

	userTests := utils.TransformValues(testData.Users, func(v data.User) TableEntry {
		return Entry(fmt.Sprint(v.ID), v.ID)
	})

	var api *API

	BeforeEach(func() {
		api = NewAPI(testData)
	})

	Context("Albums", func() {
		It("Invalid ID", func() {
			// Query
			query := `{album(id:-1){id,userid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := getData[data.Album](r, "album")

			Expect(album).To(Equal(data.Album{}))
		})

		DescribeTable("Get album by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{album(id:%v){id,userid,description}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := getData[data.Album](r, "album")

			Expect(album).To(Equal(testData.Albums[id]))
		}, albumTests)

		DescribeTable("Get album by userID", func(userId int) {
			// Query
			query := fmt.Sprintf(`{albums(userid:%v){id,userid,description}}`, userId)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			expected := utils.ValuesWhere(testData.Albums, func(album data.Album) bool { return album.UserID == userId })

			Expect(albums).To(ContainElements(expected))
		}, userTests)

		It("Get all albums", func() {
			// Query
			query := `{albums{id,userid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(Equal(testData.GetAlbums()))
		})

		Context("Album Photos", func() {
			DescribeTable("Get all album photos", func(id int) {
				query := fmt.Sprintf(`{album(id:%v){id,photos{id,albumid,description}}}`, id)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				expected := utils.ValuesWhere(testData.Photos, func(photo data.Photo) bool { return photo.AlbumID == id })

				Expect(album.Photos).To(ContainElements(expected))
			}, userTests)

			It("Get limited album photos less than size", func() {
				limit := len(testData.GetPhotosByAlbumID(0)) - 1
				// Query
				query := fmt.Sprintf(`{album(id:0){id,photos(limit:%v){id,albumid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album.Photos).To(HaveLen(limit))
			})

			It("Get limited album photos greater than size", func() {
				photos := testData.GetPhotosByAlbumID(0)
				limit := len(photos) + 1
				// Query
				query := fmt.Sprintf(`{album(id:0){id,photos(limit:%v){id,albumid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album.Photos).To(HaveLen(len(photos)))
			})
		})

		It("Get limited albums less than size", func() {
			limit := len(testData.Albums) - 1
			// Query
			query := fmt.Sprintf(`{albums(limit:%v){id,userid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(HaveLen(limit))
		})

		It("Get limited albums greater than size", func() {
			limit := len(testData.Albums) + 1
			// Query
			query := fmt.Sprintf(`{albums(limit:%v){id,userid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(HaveLen(len(testData.Albums)))
		})
	})

	Context("Photos", func() {
		It("Invalid ID", func() {
			// Query
			query := `{photo(id:-1){id,albumid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photo := getData[data.Photo](r, "photo")

			Expect(photo).To(Equal(data.Photo{}))
		})

		DescribeTable("Get photo by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{photo(id:%v){id,albumid,description}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photo := getData[data.Photo](r, "photo")

			Expect(photo).To(Equal(testData.Photos[id]))
		}, photoTests)

		DescribeTable("Get photo by albumID", func(albumId int) {
			// Query
			query := fmt.Sprintf(`{photos(albumid:%v){id,albumid,description}}`, albumId)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			expected := utils.ValuesWhere(testData.Photos, func(photo data.Photo) bool { return photo.AlbumID == albumId })

			Expect(photos).To(ContainElements(expected))
		}, userTests)

		It("Get all Photos", func() {
			// Query
			query := `{photos{id,albumid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(Equal(testData.GetPhotos()))
		})

		It("Get limited Photos less than length", func() {
			limit := len(testData.Photos) - 1
			// Query
			query := fmt.Sprintf(`{photos(limit:%v){id,albumid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(HaveLen(limit))
		})

		It("Get limited Photos greater than length", func() {
			limit := len(testData.Photos) + 1
			// Query
			query := fmt.Sprintf(`{photos(limit:%v){id,albumid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(HaveLen(len(testData.Photos)))
		})
	})

	Context("Users", func() {
		It("Invalid ID", func() {
			// Query
			query := `{user(id:-1){id,name,username}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := getData[data.User](r, "user")

			Expect(user).To(Equal(data.User{}))
		})

		DescribeTable("Get user by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{user(id:%v){id,name,username}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := getData[data.User](r, "user")

			Expect(user).To(Equal(testData.Users[id]))
		}, userTests)

		It("Get all users", func() {
			// Query
			query := `{users{id,name,username}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := getData[[]data.User](r, "users")

			Expect(users).To(Equal(testData.GetUsers()))
		})

		Context("user albums", func() {
			DescribeTable("Get user albums", func(id int) {
				query := fmt.Sprintf(`{user(id:%v){id,albums{id,userid,description}}}`, id)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				user := getData[data.User](r, "user")

				expected := utils.ValuesWhere(testData.Albums, func(album data.Album) bool { return album.UserID == id })

				Expect(user.Albums).To(ContainElements(expected))
			}, userTests)

			It("Get limited user albums less than size", func() {
				limit := len(testData.GetAlbumsByUserID(0)) - 1
				// Query
				query := fmt.Sprintf(`{user(id:0){id,albums(limit:%v){id,userid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				user := getData[data.User](r, "user")

				Expect(user.Albums).To(HaveLen(limit))
			})

			It("Get limited user albums greater than size", func() {
				albums := testData.GetAlbumsByUserID(0)
				limit := len(albums) + 1
				// Query
				query := fmt.Sprintf(`{user(id:0){id,albums(limit:%v){id,userid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				user := getData[data.User](r, "user")

				Expect(user.Albums).To(HaveLen(len(albums)))
			})
		})
		It("Get limited users less than length", func() {
			limit := len(testData.Users) - 1
			// Query
			query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := getData[[]data.User](r, "users")

			Expect(users).To(HaveLen(limit))
		})

		It("Get limited Photos greater than length", func() {
			limit := len(testData.Users) + 1
			// Query
			query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := getData[[]data.User](r, "users")

			Expect(users).To(HaveLen(len(testData.Users)))
		})
	})

	Context("Bad Schema returns err when resolving fields", func() {
		BeforeEach(func() {
			queryFields := api.Schema.QueryType().Fields()

			for _, field := range queryFields {
				// Set all fields to return the wrong type
				field.Resolve = func(p graphql.ResolveParams) (interface{}, error) { return new(interface{}), nil }
			}
		})

		It("User fields", func() {
			queries := utils.TransformValues(api.UserType.Fields(), convertFieldDefinitionToQueryString)
			for _, query := range queries {
				// Query
				query := fmt.Sprintf(`{user(id:0){%s}}`, query)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not of type data.User"))
			}
		})

		It("Photo fields", func() {
			queries := utils.TransformValues(api.PhotoType.Fields(), convertFieldDefinitionToQueryString)
			for _, query := range queries {
				// Query
				query := fmt.Sprintf(`{photo(id:0){%s}}`, query)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not of type data.Photo"))
			}
		})

		It("Album fields", func() {
			queries := utils.TransformValues(api.AlbumType.Fields(), convertFieldDefinitionToQueryString)
			for _, query := range queries {
				// Query
				query := fmt.Sprintf(`{album(id:0){%s}}`, query)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not of type data.Album"))
			}
		})

	})
})
