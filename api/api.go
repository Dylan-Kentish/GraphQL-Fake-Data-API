package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/data"
	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = "This is very secret!"

type API struct {
	Schema    graphql.Schema
	UserType  *graphql.Object
	AlbumType *graphql.Object
	PhotoType *graphql.Object
}

func NewAPI(dataModel data.IData) *API {
	photoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Photo",
		Description: "A photo.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the photo.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(photo data.Photo) interface{} {
						return photo.ID
					})
				},
			},
			"albumid": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(photo data.Photo) interface{} {
						return photo.AlbumID
					})
				},
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "The description of the photo.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(photo data.Photo) interface{} {
						return photo.Description
					})
				},
			},
		},
	})

	albumType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Album",
		Description: "A album.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(album data.Album) interface{} {
						return album.ID
					})
				},
			},
			"userid": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(album data.Album) interface{} {
						return album.UserID
					})
				},
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "The description of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(album data.Album) interface{} {
						return album.Description
					})
				},
			},
			"photos": &graphql.Field{
				Type:        graphql.NewList(photoType),
				Description: "The albums photos.",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Description: "limit the number of users",
						Type:        graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(album data.Album) interface{} {
						return utils.TryLimitIfPresent(dataModel.GetPhotosByAlbumID(album.ID), p.Args)
					})
				},
			},
		},
	})

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A user.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return user.ID
					})
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return user.Name
					})
				},
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "The username of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return user.Username
					})
				},
			},
			"email": &graphql.Field{
				Type:        graphql.String,
				Description: "The email of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return user.Email
					})
				},
			},
			"passwordHash": &graphql.Field{
				Type:        graphql.String,
				Description: "The password hash of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return string(user.PasswordHash[:])
					})
				},
			},
			"albums": &graphql.Field{
				Type:        graphql.NewList(albumType),
				Description: "The users albums.",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Description: "limit the number of users",
						Type:        graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return utils.TryLimitIfPresent(dataModel.GetAlbumsByUserID(user.ID), p.Args)
					})
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        userType,
				Description: "User by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the user",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dataModel.GetUser(p.Args["id"].(int)), nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "All users",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Description: "limit the number of users",
						Type:        graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return utils.TryLimitIfPresent(dataModel.GetUsers(), p.Args), nil
				},
			},
			"album": &graphql.Field{
				Type:        albumType,
				Description: "Album by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the album",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dataModel.GetAlbum(p.Args["id"].(int)), nil
				},
			},
			"albums": &graphql.Field{
				Type:        graphql.NewList(albumType),
				Description: "All albums",
				Args: graphql.FieldConfigArgument{
					"userid": &graphql.ArgumentConfig{
						Description: "id of the user",
						Type:        graphql.Int,
					},
					"limit": &graphql.ArgumentConfig{
						Description: "limit the number of albums",
						Type:        graphql.Int,
					},
				},

				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var albums []data.Album

					if id, exists := p.Args["userid"].(int); exists {
						albums = dataModel.GetAlbumsByUserID(id)
					} else {
						albums = dataModel.GetAlbums()
					}

					return utils.TryLimitIfPresent(albums, p.Args), nil
				},
			},
			"photo": &graphql.Field{
				Type:        photoType,
				Description: "Photo by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the photo",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dataModel.GetPhoto(p.Args["id"].(int)), nil
				},
			},
			"photos": &graphql.Field{
				Type:        graphql.NewList(photoType),
				Description: "All photos",
				Args: graphql.FieldConfigArgument{
					"albumid": &graphql.ArgumentConfig{
						Description: "id of the album",
						Type:        graphql.Int,
					},
					"limit": &graphql.ArgumentConfig{
						Description: "limit the number of photos",
						Type:        graphql.Int,
					},
				},

				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var photos []data.Photo

					if id, exists := p.Args["albumid"].(int); exists {
						photos = dataModel.GetPhotosByAlbumID(id)
					} else {
						photos = dataModel.GetPhotos()
					}

					return utils.TryLimitIfPresent(photos, p.Args), nil
				},
			},
		},
	})

	authenticationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Authentication",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type:        graphql.String,
				Description: "Authentication token",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(auth data.Authentication) interface{} {
						return auth.Token
					})
				},
			},
			"user": &graphql.Field{
				Type:        userType,
				Description: "User",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(auth data.Authentication) interface{} {
						return auth.User
					})
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type:        authenticationType,
				Description: "User authentication",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Description: "email of the user",
						Type:        graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Description: "password of the user",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := dataModel.GetUserWithEmail(p.Args["email"].(string))

					if err != nil {
						return nil, errors.New("invalid email or password")
					}

					if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p.Args["password"].(string))); err != nil {
						return nil, errors.New("invalid email or password")
					}

					claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
						Issuer:    strconv.Itoa(user.ID),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					})

					token, err := claims.SignedString([]byte(secretKey))

					if err != nil {
						// What would cause this to happen?
						return nil, errors.New("login failed")
					}

					return data.Authentication{
						Token: token,
						User:  *user,
					}, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})

	return &API{
		Schema:    schema,
		UserType:  userType,
		AlbumType: albumType,
		PhotoType: photoType,
	}
}

func resolveType[T any](p graphql.ResolveParams, onOk func(model T) interface{}) (interface{}, error) {
	if model, ok := p.Source.(T); ok {
		return onOk(model), nil
	}
	return nil, fmt.Errorf("source is not of type %T", *new(T))
}
