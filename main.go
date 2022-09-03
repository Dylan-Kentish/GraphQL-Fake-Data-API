package main

import (
	"fmt"
	"net/http"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/api"
	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/testUtils"
	"github.com/graphql-go/handler"
)

func main() {
	data := testUtils.NewTestData()
	api := api.NewAPI(data)

	h := handler.New(&handler.Config{
		Schema:   &api.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	fmt.Println("Starting server at localhost:8080/graphql")
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
