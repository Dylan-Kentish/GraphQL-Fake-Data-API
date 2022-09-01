package main

import (
	"fmt"
	"net/http"

	api "github.com/Dylan-Kentish/GraphQLFakeDataAPI/internal"
	"github.com/graphql-go/handler"
)

func main() {
	h := handler.New(&handler.Config{
		Schema:   &api.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	fmt.Println("Starting server at localhost:8080/graphql")
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
