package api

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"golang.org/x/exp/maps"
)

var (
	UserData   map[int]User
	UserSchema graphql.Schema
)

type User struct {
	ID       string
	Name     string
	Username string
}

func init() {
	UserData = GetUserData()

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A user.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The id of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(User); ok {
						return human.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(User); ok {
						return human.Name, nil
					}
					return nil, nil
				},
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "The username of the user.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(User); ok {
						return human.Username, nil
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
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, err := strconv.Atoi(p.Args["id"].(string))
					if err != nil {
						return nil, err
					}
					return GetUser(id), nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "All users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return maps.Values(UserData), nil
				},
			},
		},
	})

	UserSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}

func GetUser(id int) User {
	if user, ok := UserData[id]; ok {
		return user
	}
	return User{}
}
