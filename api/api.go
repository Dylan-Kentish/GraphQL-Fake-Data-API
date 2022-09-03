package api

import (
	"errors"

	"github.com/graphql-go/graphql"
	"golang.org/x/exp/maps"
)

type API struct {
	Schema    graphql.Schema
	UserType  *graphql.Object
	AlbumType *graphql.Object
	PhotoType *graphql.Object
}

func NewAPI(data *data) *API {
	photoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Photo",
		Description: "A photo.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the photo.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if album, ok := p.Source.(Photo); ok {
						return album.ID, nil
					}
					return nil, errors.New("source is not a api.Photo")
				},
			},
			"albumid": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if album, ok := p.Source.(Photo); ok {
						return album.AlbumID, nil
					}
					return nil, errors.New("source is not a api.Photo")
				},
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "The description of the photo.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if album, ok := p.Source.(Photo); ok {
						return album.Description, nil
					}
					return nil, errors.New("source is not a api.Photo")
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
					if album, ok := p.Source.(Album); ok {
						return album.ID, nil
					}
					return nil, errors.New("source is not a api.Album")
				},
			},
			"userid": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if album, ok := p.Source.(Album); ok {
						return album.UserID, nil
					}
					return nil, errors.New("source is not a api.Album")
				},
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "The description of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if album, ok := p.Source.(Album); ok {
						return album.Description, nil
					}
					return nil, errors.New("source is not a api.Album")
				},
			},
			"photos": &graphql.Field{
				Type:        graphql.NewList(photoType),
				Description: "The albums photos.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(Album); ok {
						return data.getPhotosByAlbumID(user.ID), nil
					}
					return nil, errors.New("source is not a api.Album")
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
					if user, ok := p.Source.(User); ok {
						return user.ID, nil
					}
					return nil, errors.New("source is not a api.User")
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.Name, nil
					}
					return nil, errors.New("source is not a api.User")
				},
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "The username of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.Username, nil
					}
					return nil, errors.New("source is not a api.User")
				},
			},
			"albums": &graphql.Field{
				Type:        graphql.NewList(albumType),
				Description: "The users albums.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return data.getAlbumsByUserID(user.ID), nil
					}
					return nil, errors.New("source is not a api.User")
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
					return data.getUser(p.Args["id"].(int)), nil
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
					albums := maps.Values(data.Users)
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
					return data.getAlbum(p.Args["id"].(int)), nil
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
					var albums []Album

					if id, exists := p.Args["userid"].(int); exists {
						albums = data.getAlbumsByUserID(id)
					} else {
						albums = maps.Values(data.Albums)
					}

					if limit, exists := p.Args["limit"].(int); exists {
						return albums[:limit], nil
					} else {
						return albums, nil
					}
				},
			},
			"photo": &graphql.Field{
				Type:        albumType,
				Description: "Photo by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the photo",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return data.getPhoto(p.Args["id"].(int)), nil
				},
			},
			"photos": &graphql.Field{
				Type:        graphql.NewList(albumType),
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
					var photos []Photo

					if id, exists := p.Args["albumid"].(int); exists {
						photos = data.getPhotosByAlbumID(id)
					} else {
						photos = maps.Values(data.Photos)
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
