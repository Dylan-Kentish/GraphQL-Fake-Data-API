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

	photoTests := utils.Transform(testData.GetPhotos(), func(v data.Photo) TableEntry {
		return Entry(fmt.Sprint(v.ID), v.ID)
	})

	albumTests := utils.Transform(testData.GetAlbums(), func(v data.Album) TableEntry {
		return Entry(fmt.Sprint(v.ID), v.ID)
	})

	userTests := utils.Transform(testData.GetUsers(), func(v data.User) TableEntry {
		return Entry(fmt.Sprint(v.ID), v.ID)
	})

	var api *API

	BeforeEach(func() {
		api = NewAPI(testData)
	})

	Context("Albums", func() {
		It("Invalid ID", func() {

			query := `{album(id:-1){id,userid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := getData[data.Album](r, "album")

			Expect(album).To(Equal(data.Album{}))
		})

		DescribeTable("Get album by ID", func(id int) {
			expected := testData.GetAlbum(id)

			query := fmt.Sprintf(`{album(id:%v){id,userid,description}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := getData[data.Album](r, "album")

			Expect(album).To(Equal(expected))
		}, albumTests)

		DescribeTable("Get album by userID", func(id int) {
			expected := testData.GetAlbumsByUserID(id)

			query := fmt.Sprintf(`{albums(userid:%v){id,userid,description}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(Equal(expected))
		}, userTests)

		It("Get all albums", func() {
			expected := testData.GetAlbums()

			query := `{albums{id,userid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(Equal(expected))
		})

		Context("Album Photos", func() {
			DescribeTable("Get all album photos", func(id int) {
				expected := testData.GetPhotosByAlbumID(id)

				query := fmt.Sprintf(`{album(id:%v){id,photos{id,albumid,description}}}`, id)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album.Photos).To(Equal(expected))
			}, userTests)

			It("Get limited album photos less than size", func() {
				limit := len(testData.GetPhotosByAlbumID(0)) - 1

				query := fmt.Sprintf(`{album(id:0){id,photos(limit:%v){id,albumid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album.Photos).To(HaveLen(limit))
			})

			It("Get limited album photos greater than size", func() {
				expected := len(testData.GetPhotosByAlbumID(0))
				limit := expected + 1

				query := fmt.Sprintf(`{album(id:0){id,photos(limit:%v){id,albumid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album.Photos).To(HaveLen(expected))
			})
		})

		It("Get limited albums less than size", func() {
			limit := len(testData.GetAlbums()) - 1

			query := fmt.Sprintf(`{albums(limit:%v){id,userid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(HaveLen(limit))
		})

		It("Get limited albums greater than size", func() {
			expected := len(testData.GetAlbums())
			limit := expected + 1

			query := fmt.Sprintf(`{albums(limit:%v){id,userid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := getData[[]data.Album](r, "albums")

			Expect(albums).To(HaveLen(expected))
		})
	})

	Context("Photos", func() {
		It("Invalid ID", func() {

			query := `{photo(id:-1){id,albumid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photo := getData[data.Photo](r, "photo")

			Expect(photo).To(Equal(data.Photo{}))
		})

		DescribeTable("Get photo by ID", func(id int) {
			expected := testData.GetPhoto(id)

			query := fmt.Sprintf(`{photo(id:%v){id,albumid,description}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photo := getData[data.Photo](r, "photo")

			Expect(photo).To(Equal(expected))
		}, photoTests)

		DescribeTable("Get photo by albumID", func(id int) {
			expected := testData.GetPhotosByAlbumID(id)

			query := fmt.Sprintf(`{photos(albumid:%v){id,albumid,description}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(ContainElements(expected))
		}, userTests)

		It("Get all Photos", func() {
			expected := testData.GetPhotos()

			query := `{photos{id,albumid,description}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(Equal(expected))
		})

		It("Get limited Photos less than length", func() {
			limit := len(testData.GetPhotos()) - 1

			query := fmt.Sprintf(`{photos(limit:%v){id,albumid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(HaveLen(limit))
		})

		It("Get limited Photos greater than length", func() {
			expected := len(testData.GetPhotos())
			limit := expected + 1

			query := fmt.Sprintf(`{photos(limit:%v){id,albumid,description}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := getData[[]data.Photo](r, "photos")

			Expect(photos).To(HaveLen(expected))
		})
	})

	Context("Users", func() {
		It("Invalid ID", func() {
			query := `{user(id:-1){id,name,username}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := getData[data.User](r, "user")

			Expect(user).To(Equal(data.User{}))
		})

		DescribeTable("Get user by ID", func(id int) {
			expected := testData.GetUser(id)

			query := fmt.Sprintf(`{user(id:%v){id,name,username}}`, id)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := getData[data.User](r, "user")

			Expect(user).To(Equal(expected))
		}, userTests)

		It("Get all users", func() {
			expected := testData.GetUsers()

			query := `{users{id,name,username}}`
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := getData[[]data.User](r, "users")

			Expect(users).To(Equal(expected))
		})

		Context("user albums", func() {
			DescribeTable("Get user albums", func(id int) {
				expected := testData.GetAlbumsByUserID(id)

				query := fmt.Sprintf(`{user(id:%v){id,albums{id,userid,description}}}`, id)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				user := getData[data.User](r, "user")

				Expect(user.Albums).To(Equal(expected))
			}, userTests)

			It("Get limited user albums less than size", func() {
				limit := len(testData.GetAlbumsByUserID(0)) - 1

				query := fmt.Sprintf(`{user(id:0){id,albums(limit:%v){id,userid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				user := getData[data.User](r, "user")

				Expect(user.Albums).To(HaveLen(limit))
			})

			It("Get limited user albums greater than size", func() {
				expected := len(testData.GetAlbumsByUserID(0))
				limit := expected + 1

				query := fmt.Sprintf(`{user(id:0){id,albums(limit:%v){id,userid,description}}}`, limit)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				user := getData[data.User](r, "user")

				Expect(user.Albums).To(HaveLen(expected))
			})
		})
		It("Get limited users less than length", func() {
			limit := len(testData.GetUsers()) - 1

			query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := getData[[]data.User](r, "users")

			Expect(users).To(HaveLen(limit))
		})

		It("Get limited Photos greater than length", func() {
			expected := len(testData.GetUsers())
			limit := expected + 1

			query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
			params := graphql.Params{Schema: api.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := getData[[]data.User](r, "users")

			Expect(users).To(HaveLen(expected))
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
				query := fmt.Sprintf(`{album(id:0){%s}}`, query)
				params := graphql.Params{Schema: api.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not of type data.Album"))
			}
		})

	})
})
