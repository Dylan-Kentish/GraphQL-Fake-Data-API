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
		var variables map[string]interface{}
		var params graphql.Params
		var query string

		BeforeEach(func() {
			variables = make(map[string]interface{})
			params = graphql.Params{
				Schema:         api.Schema,
				VariableValues: variables,
			}
		})

		JustBeforeEach(func() {
			params.RequestString = query
		})

		Context("album", func() {
			BeforeEach(func() {
				query = `
					query ($id: Int!, $withPhotos: Boolean = false, $limit: Int) {
						album(id:$id){
							id
							userid
							description
							photos(limit:$limit) @include(if: $withPhotos) {
								id
								albumid
								description
							}
						}
					}`
				variables["limit"] = nil
			})

			It("Invalid ID", func() {
				variables["id"] = -1

				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album).To(Equal(data.Album{}))
			})

			DescribeTable("Get album by ID", func(id int) {
				variables["id"] = id
				expected := testData.GetAlbum(id)

				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album).To(Equal(expected))
			}, albumTests)

			DescribeTable("Get all album photos", func(id int) {
				variables["id"] = id
				variables["withPhotos"] = true
				expected := testData.GetPhotosByAlbumID(id)

				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				album := getData[data.Album](r, "album")

				Expect(album.Photos).To(Equal(expected))
			}, userTests)

			Context("limited album photos", func() {
				var limit int
				var albumId = 0
				JustBeforeEach(func() {
					variables["id"] = albumId
					variables["withPhotos"] = true
					variables["limit"] = limit
				})

				When("limited less than size", func() {
					BeforeEach(func() {
						limit = len(testData.GetPhotosByAlbumID(albumId)) - 1
					})

					It("returns limit", func() {
						r := graphql.Do(params)
						Expect(r.Errors).To(BeEmpty())

						album := getData[data.Album](r, "album")

						Expect(album.Photos).To(HaveLen(limit))
					})
				})

				When("limited greter than size", func() {
					BeforeEach(func() {
						limit = len(testData.GetPhotosByAlbumID(albumId)) + 1
					})

					It("returns size", func() {
						expected := len(testData.GetPhotosByAlbumID(albumId))

						r := graphql.Do(params)
						Expect(r.Errors).To(BeEmpty())

						album := getData[data.Album](r, "album")

						Expect(album.Photos).To(HaveLen(expected))
					})
				})
			})
		})

		Context("albums", func() {
			BeforeEach(func() {
				query = `
				query {
					albums {
						id
						userid
						description
					}
				}`
			})

			It("Get all albums", func() {
				expected := testData.GetAlbums()

				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				albums := getData[[]data.Album](r, "albums")

				Expect(albums).To(Equal(expected))
			})

			Context("albums by userID", func() {
				BeforeEach(func() {
					query = `
					query ($id: Int!) {
						albums(userid: $id) {
							id
							userid
							description
						}
					}`
				})

				DescribeTable("Get albums by userID", func(id int) {
					variables["id"] = id
					expected := testData.GetAlbumsByUserID(id)

					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					albums := getData[[]data.Album](r, "albums")

					Expect(albums).To(Equal(expected))
				}, userTests)
			})

			Context("limited albums", func() {
				var limit int

				BeforeEach(func() {
					query = `
						query ($limit: Int!) {
							albums(limit:$limit) {
								id
								userid
								description
							}
						}`
				})

				JustBeforeEach(func() {
					variables["limit"] = limit
				})

				When("limited less than size", func() {
					BeforeEach(func() {
						limit = len(testData.GetAlbums()) - 1
					})

					It("returns limit", func() {
						r := graphql.Do(params)
						Expect(r.Errors).To(BeEmpty())

						albums := getData[[]data.Album](r, "albums")

						Expect(albums).To(HaveLen(limit))
					})
				})

				When("limited greter than size", func() {
					BeforeEach(func() {
						limit = len(testData.GetAlbums()) + 1
					})

					It("returns size", func() {
						expected := len(testData.GetAlbums())

						r := graphql.Do(params)
						Expect(r.Errors).To(BeEmpty())

						albums := getData[[]data.Album](r, "albums")

						Expect(albums).To(HaveLen(expected))
					})
				})
			})
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
