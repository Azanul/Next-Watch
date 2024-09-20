package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Azanul/Next-Watch/graph"
	"github.com/Azanul/Next-Watch/internal/auth"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{
			Resolvers: &graph.Resolver{},
			Directives: graph.DirectiveRoot{
				HasRole: hasRoleDirective,
			},
		},
	))

	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func hasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role string) (interface{}, error) {
	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the user has the required role
	if user.Role != role {
		return nil, errors.New("access denied")
	}

	// If the user has the role, continue to the next resolver
	return next(ctx)
}

//func cors(h http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "*")
//		w.Header().Set("Access-Control-Allow-Headers", "*")
//		h(w, r)
//	}
//}
