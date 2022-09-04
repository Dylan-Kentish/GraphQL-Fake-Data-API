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

	var api *API

	Context("Valid Schema", func() {
		photoTests := utils.Transform(testData.GetPhotos(), func(v data.Photo) TableEntry {
			return Entry(fmt.Sprint(v.ID), v.ID)
		})

		albumTests := utils.Transform(testData.GetAlbums(), func(v data.Album) TableEntry {
			return Entry(fmt.Sprint(v.ID), v.ID)
		})

		userTests := utils.Transform(testData.GetUsers(), func(v data.User) TableEntry {
			return Entry(fmt.Sprint(v.ID), v.ID)
		})

		var variables map[string]interface{}
		var params graphql.Params
		var query string

		BeforeEach(func() {
			api = NewAPI(testData)
			variables = make(map[string]interface{})
			params = graphql.Params{
				Schema:         api.Schema,
				VariableValues: variables,
			}
		})

		JustBeforeEach(func() {
			params.RequestString = query
		})

		Context("Albums", func() {
			Context("album", func() {
				BeforeEach(func() {
					query = `
						query ($id: Int!, $withPhotos: Boolean = false) {
							album(id:$id){
								id
								userid
								description
								photos @include(if: $withPhotos) {
									id
									albumid
									description
								}
							}
						}`
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

					BeforeEach(func() {
						query = `
							query ($id: Int!, $limit: Int) {
								album(id:$id){
									id
									userid
									description
									photos(limit:$limit) {
										id
										albumid
										description
									}
								}
							}`
					})

					JustBeforeEach(func() {
						variables["id"] = albumId
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
			Context("photo", func() {
				BeforeEach(func() {
					query = `
						query ($id: Int!) {
							photo(id:$id) {
								id 
								albumid
								description
							}
						}`
				})

				It("Invalid ID", func() {
					variables["id"] = -1
					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					photo := getData[data.Photo](r, "photo")

					Expect(photo).To(Equal(data.Photo{}))
				})

				DescribeTable("Get photo by ID", func(id int) {
					variables["id"] = id
					expected := testData.GetPhoto(id)

					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					photo := getData[data.Photo](r, "photo")

					Expect(photo).To(Equal(expected))
				}, photoTests)
			})

			Context("photos", func() {
				BeforeEach(func() {
					query = `
						query {
							photos {
								id 
								albumid
								description
							}
						}`
				})

				It("Get all Photos", func() {
					expected := testData.GetPhotos()

					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					photos := getData[[]data.Photo](r, "photos")

					Expect(photos).To(Equal(expected))
				})

				Context("photos by albumID", func() {
					BeforeEach(func() {
						query = `
						query ($id: Int!) {
							photos(albumid: $id) {
								id 
								albumid
								description
							}
						}`
					})

					DescribeTable("Get photos by albumID", func(id int) {
						variables["id"] = id
						expected := testData.GetPhotosByAlbumID(id)

						r := graphql.Do(params)
						Expect(r.Errors).To(BeEmpty())

						photos := getData[[]data.Photo](r, "photos")

						Expect(photos).To(ContainElements(expected))
					}, userTests)
				})

				Context("limited photos", func() {
					var limit int

					BeforeEach(func() {
						query = `
							query ($limit: Int) {
								photos(limit:$limit) {
									id
									albumid
									description
								}
							}`
					})

					JustBeforeEach(func() {
						variables["limit"] = limit
					})

					When("limited less than size", func() {
						BeforeEach(func() {
							limit = len(testData.GetPhotos()) - 1
						})

						It("returns limit", func() {
							r := graphql.Do(params)
							Expect(r.Errors).To(BeEmpty())

							photos := getData[[]data.Photo](r, "photos")

							Expect(photos).To(HaveLen(limit))
						})
					})

					When("limited greter than size", func() {
						BeforeEach(func() {
							limit = len(testData.GetPhotos()) + 1
						})

						It("returns size", func() {
							expected := len(testData.GetPhotos())

							r := graphql.Do(params)
							Expect(r.Errors).To(BeEmpty())

							photos := getData[[]data.Photo](r, "photos")

							Expect(photos).To(HaveLen(expected))
						})
					})
				})
			})
		})

		Context("Users", func() {
			Context("user", func() {
				BeforeEach(func() {
					query = `
						query ($id: Int!, $withPhotos: Boolean = false) {
							user(id:$id){
								id
								name
								username
								email
								passwordHash
								albums @include(if: $withPhotos) {
									id
									userid
									description
								}
							}
						}`
				})

				It("Invalid ID", func() {
					variables["id"] = -1
					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					user := getData[data.User](r, "user")

					Expect(user).To(Equal(data.User{}))
				})

				DescribeTable("Get user by ID", func(id int) {
					variables["id"] = id
					expected := testData.GetUser(id)

					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					user := getData[data.User](r, "user")

					Expect(user).To(Equal(expected))
				}, userTests)

				DescribeTable("Get user albums", func(id int) {
					variables["withPhotos"] = true
					variables["id"] = id
					expected := testData.GetAlbumsByUserID(id)

					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					user := getData[data.User](r, "user")

					Expect(user.Albums).To(Equal(expected))
				}, userTests)

				Context("limited user albums", func() {
					var limit int

					BeforeEach(func() {
						query = `
							query ($limit: Int) {
								user(id:0) {
									id
									albums(limit:$limit) {
										id
										userid
										description
									}
								}
							}`
					})

					JustBeforeEach(func() {
						variables["limit"] = limit
					})

					When("limited less than size", func() {
						BeforeEach(func() {
							limit = len(testData.GetAlbumsByUserID(0)) - 1
						})

						It("returns limit", func() {
							r := graphql.Do(params)
							Expect(r.Errors).To(BeEmpty())

							user := getData[data.User](r, "user")

							Expect(user.Albums).To(HaveLen(limit))
						})
					})

					When("limited greter than size", func() {
						BeforeEach(func() {
							limit = len(testData.GetAlbumsByUserID(0)) + 1
						})

						It("returns size", func() {
							expected := len(testData.GetAlbumsByUserID(0))

							r := graphql.Do(params)
							Expect(r.Errors).To(BeEmpty())

							user := getData[data.User](r, "user")

							Expect(user.Albums).To(HaveLen(expected))
						})
					})
				})
			})

			Context("users", func() {
				BeforeEach(func() {
					query = `
						query {
							users {
								id
								name
								username
								email
								passwordHash
							}
						}`
				})

				It("Get all users", func() {
					expected := testData.GetUsers()

					r := graphql.Do(params)
					Expect(r.Errors).To(BeEmpty())

					users := getData[[]data.User](r, "users")

					Expect(users).To(Equal(expected))
				})

				Context("limited users", func() {
					var limit int

					BeforeEach(func() {
						query = `
						query ($limit: Int) {
							users(limit:$limit) {
								id
								name
								username
							}
							}`
					})

					JustBeforeEach(func() {
						variables["limit"] = limit
					})

					When("limited less than size", func() {
						BeforeEach(func() {
							limit = len(testData.GetUsers()) - 1
						})

						It("returns limit", func() {
							r := graphql.Do(params)
							Expect(r.Errors).To(BeEmpty())

							users := getData[[]data.User](r, "users")

							Expect(users).To(HaveLen(limit))
						})
					})

					When("limited greter than size", func() {
						BeforeEach(func() {
							limit = len(testData.GetUsers()) + 1
						})

						It("returns size", func() {
							expected := len(testData.GetUsers())

							r := graphql.Do(params)
							Expect(r.Errors).To(BeEmpty())

							users := getData[[]data.User](r, "users")

							Expect(users).To(HaveLen(expected))
						})
					})
				})
			})
		})
	})

	Context("Bad Schema", func() {
		BeforeEach(func() {
			queryFields := api.Schema.QueryType().Fields()

			for _, field := range queryFields {
				// Set all fields to return the wrong type
				field.Resolve = func(p graphql.ResolveParams) (interface{}, error) { return new(interface{}), nil }
			}
		})

		When("resolving User fields", func() {
			It("returns err", func() {
				queries := utils.TransformValues(api.UserType.Fields(), convertFieldDefinitionToQueryString)
				for _, query := range queries {
					query := fmt.Sprintf(`{user(id:0){%s}}`, query)
					params := graphql.Params{Schema: api.Schema, RequestString: query}
					r := graphql.Do(params)
					Expect(r.Errors).To(HaveLen(1))
					Expect(r.Errors[0].Message).To(Equal("source is not of type data.User"))
				}
			})
		})

		When("resolving Photo fields", func() {
			It("returns err", func() {
				queries := utils.TransformValues(api.PhotoType.Fields(), convertFieldDefinitionToQueryString)
				for _, query := range queries {
					query := fmt.Sprintf(`{photo(id:0){%s}}`, query)
					params := graphql.Params{Schema: api.Schema, RequestString: query}
					r := graphql.Do(params)
					Expect(r.Errors).To(HaveLen(1))
					Expect(r.Errors[0].Message).To(Equal("source is not of type data.Photo"))
				}
			})
		})

		When("resolving Album fields", func() {
			It("returns err", func() {
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

	Context("Authentication", func() {
		var variables map[string]interface{}
		var params graphql.Params
		var mutation string

		BeforeEach(func() {
			api = NewAPI(testData)
			variables = make(map[string]interface{})
			params = graphql.Params{
				Schema:         api.Schema,
				VariableValues: variables,
			}
		})

		JustBeforeEach(func() {
			params.RequestString = mutation
		})

		When("login", func() {
			BeforeEach(func() {
				mutation = `
					mutation ($email: String!, $password: String!) {
						login(email:$email, password:$password){
							token
							user {
								id
							}
						}
					}`
			})

			It("Authenticates valid data", func() {
				user := testData.GetUser(0)
				variables["email"] = user.Email
				variables["password"] = "Password0"

				r := graphql.Do(params)
				Expect(r.Errors).To(BeEmpty())

				authentication := getData[data.Authentication](r, "authentication")

				Expect(authentication.Token).ToNot(BeNil())
				Expect(authentication.User.ID).To(Equal(user.ID))
			})

		})

		When("Invalid Credentials", func() {
			It("Rejects invalid email", func() {
				variables["email"] = "not their email"
				variables["password"] = "Password0"

				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("invalid email or password"))
			})

			It("Rejects invalid password", func() {
				user := testData.GetUser(0)
				variables["email"] = user.Email
				variables["password"] = "not their password"

				r := graphql.Do(params)
				Expect(r.Errors).To(HaveLen(1))
				Expect(r.Errors[0].Message).To(Equal("invalid email or password"))
			})
		})
	})
})
