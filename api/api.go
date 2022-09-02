package api

import (
	"github.com/graphql-go/graphql"
	"golang.org/x/exp/maps"
)

var (
	Data   data
	Schema graphql.Schema
)

func init() {
	Data = generateData()

	albumType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Album",
		Description: "A album.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(Album); ok {
						return user.ID, nil
					}
					return nil, nil
				},
			},
			"userid": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The id of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(Album); ok {
						return user.UserID, nil
					}
					return nil, nil
				},
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "The description of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(Album); ok {
						return user.Description, nil
					}
					return nil, nil
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
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.Name, nil
					}
					return nil, nil
				},
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "The username of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(User); ok {
						return user.Username, nil
					}
					return nil, nil
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
					return getUser(p.Args["id"].(int)), nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "All users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return maps.Values(Data.Users), nil
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
					return getAlbum(p.Args["id"].(int)), nil
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
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, exists := p.Args["userid"].(int); exists {
						return getAlbumsByUserID(id), nil
					}

					return maps.Values(Data.Albums), nil
				},
			},
		},
	})

	Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}

func getUser(id int) User {
	if user, ok := Data.Users[id]; ok {
		return user
	}
	return User{}
}

func getAlbum(id int) Album {
	if album, ok := Data.Albums[id]; ok {
		return album
	}
	return Album{}
}

func getAlbumsByUserID(userID int) []Album {
	albums := make([]Album, 0)

	for _, album := range Data.Albums {
		if album.UserID == userID {
			albums = append(albums, album)
		}
	}

	return albums
}
