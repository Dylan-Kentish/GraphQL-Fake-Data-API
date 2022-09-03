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

	var testApi *API

	BeforeEach(func() {
		testApi = NewAPI(testData)
	})

	Context("Albums", func() {
		It("Invalid ID", func() {
			// Query
			query := `{album(id:-1){id,userid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := testUtils.GetData[data.Album](r, "album")

			Expect(album).To(Equal(data.Album{}))
		})

		DescribeTable("Get album by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{album(id:%v){id,userid,description}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := testUtils.GetData[data.Album](r, "album")

			Expect(album).To(Equal(testData.Albums[id]))
		}, albumTests)

		DescribeTable("Get album by userID", func(userId int) {
			// Query
			query := fmt.Sprintf(`{albums(userid:%v){id,userid,description}}`, userId)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := testUtils.GetData[[]data.Album](r, "albums")

			expected := utils.ValuesWhere(testData.Albums, func(album data.Album) bool { return album.UserID == userId })

			Expect(albums).To(ContainElements(expected))
		}, userTests)

		It("Get all albums", func() {
			// Query
			query := `{albums{id,userid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := testUtils.GetData[[]data.Album](r, "albums")

			Expect(albums).To(ContainElements(maps.Values(testData.Albums)))
		})

		DescribeTable("Get album photos", func(id int) {
			query := fmt.Sprintf(`{album(id:%v){id,photos{id,albumid,description}}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			album := testUtils.GetData[data.Album](r, "album")

			expected := utils.ValuesWhere(testData.Photos, func(photo data.Photo) bool { return photo.AlbumID == id })

			Expect(album.Photos).To(ContainElements(expected))
		}, userTests)

		It("Get limited albums", func() {
			limit := 5
			// Query
			query := fmt.Sprintf(`{albums(limit:%v){id,userid,description}}`, limit)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			albums := testUtils.GetData[[]data.Album](r, "albums")

			Expect(albums).To(HaveLen(limit))
		})
	})

	Context("Photos", func() {
		It("Invalid ID", func() {
			// Query
			query := `{photo(id:-1){id,albumid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photo := testUtils.GetData[data.Photo](r, "photo")

			Expect(photo).To(Equal(data.Photo{}))
		})

		DescribeTable("Get photo by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{photo(id:%v){id,albumid,description}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photo := testUtils.GetData[data.Photo](r, "photo")

			Expect(photo).To(Equal(testData.Photos[id]))
		}, photoTests)

		DescribeTable("Get photo by albumID", func(albumId int) {
			// Query
			query := fmt.Sprintf(`{photos(albumid:%v){id,albumid,description}}`, albumId)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := testUtils.GetData[[]data.Photo](r, "photos")

			expected := utils.ValuesWhere(testData.Photos, func(photo data.Photo) bool { return photo.AlbumID == albumId })

			Expect(photos).To(ContainElements(expected))
		}, userTests)

		It("Get all Photos", func() {
			// Query
			query := `{photos{id,albumid,description}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := testUtils.GetData[[]data.Photo](r, "photos")

			Expect(photos).To(ContainElements(maps.Values(testData.Photos)))
		})

		It("Get limited Photos", func() {
			limit := 5
			// Query
			query := fmt.Sprintf(`{photos(limit:%v){id,albumid,description}}`, limit)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			photos := testUtils.GetData[[]data.Photo](r, "photos")

			Expect(photos).To(HaveLen(limit))
		})
	})

	Context("Users", func() {
		It("Invalid ID", func() {
			// Query
			query := `{user(id:-1){id,name,username}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := testUtils.GetData[data.User](r, "user")

			Expect(user).To(Equal(data.User{}))
		})

		DescribeTable("Get user by ID", func(id int) {
			// Query
			query := fmt.Sprintf(`{user(id:%v){id,name,username}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := testUtils.GetData[data.User](r, "user")

			Expect(user).To(Equal(testData.Users[id]))
		}, userTests)

		It("Get all users", func() {
			// Query
			query := `{users{id,name,username}}`
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := testUtils.GetData[[]data.User](r, "users")

			Expect(users).To(ContainElements(maps.Values(testData.Users)))
		})

		DescribeTable("Get user albums", func(id int) {
			query := fmt.Sprintf(`{user(id:%v){id,albums{id,userid,description}}}`, id)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			user := testUtils.GetData[data.User](r, "user")

			expected := utils.ValuesWhere(testData.Albums, func(album data.Album) bool { return album.UserID == id })

			Expect(user.Albums).To(ContainElements(expected))
		}, userTests)

		It("Get limited users", func() {
			limit := 5
			// Query
			query := fmt.Sprintf(`{users(limit:%v){id,name,username}}`, limit)
			params := graphql.Params{Schema: testApi.Schema, RequestString: query}
			r := graphql.Do(params)
			Expect(r.Errors).To(BeEmpty())

			users := testUtils.GetData[[]data.User](r, "users")

			Expect(users).To(HaveLen(limit))
		})
	})

	Context("Bad Schema returns err when resolving fields", func() {
		BeforeEach(func() {
			queryFields := testApi.Schema.QueryType().Fields()

			for _, field := range queryFields {
				// Set all fields to return the wrong type
				field.Resolve = func(p graphql.ResolveParams) (interface{}, error) { return new(interface{}), nil }
			}
		})

		It("User fields", func() {
			queries := utils.TransformValues(testApi.UserType.Fields(), testUtils.ConvertFieldDefinitionToQueryString)
			for _, query := range queries {
				// Query
				query := fmt.Sprintf(`{user(id:0){%s}}`, query)
				params := graphql.Params{Schema: testApi.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not a api.User"))
			}
		})

		It("Photo fields", func() {
			queries := utils.TransformValues(testApi.PhotoType.Fields(), testUtils.ConvertFieldDefinitionToQueryString)
			for _, query := range queries {
				// Query
				query := fmt.Sprintf(`{photo(id:0){%s}}`, query)
				params := graphql.Params{Schema: testApi.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not a api.Photo"))
			}
		})

		It("Album fields", func() {
			queries := utils.TransformValues(testApi.AlbumType.Fields(), testUtils.ConvertFieldDefinitionToQueryString)
			for _, query := range queries {
				// Query
				query := fmt.Sprintf(`{album(id:0){%s}}`, query)
				params := graphql.Params{Schema: testApi.Schema, RequestString: query}
				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("source is not a api.Album"))
			}
		})

	})
})
