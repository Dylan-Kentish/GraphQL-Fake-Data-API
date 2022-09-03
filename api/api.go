package api

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/data"
	"github.com/graphql-go/graphql"
)

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
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(album data.Album) interface{} {
						return dataModel.GetPhotosByAlbumID(album.ID)
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
			"albums": &graphql.Field{
				Type:        graphql.NewList(albumType),
				Description: "The users albums.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveType(p, func(user data.User) interface{} {
						return dataModel.GetAlbumsByUserID(user.ID)
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
					albums := dataModel.GetUsers()
					if limit, exists := p.Args["limit"].(int); exists {
						return albums[:limit], nil
					} else {
						return albums, nil
					}
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

					if limit, exists := p.Args["limit"].(int); exists {
						return albums[:limit], nil
					} else {
						return albums, nil
					}
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
				Description: "All albums",
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

					if limit, exists := p.Args["limit"].(int); exists {
						return photos[:limit], nil
					} else {
						return photos, nil
					}
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
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
