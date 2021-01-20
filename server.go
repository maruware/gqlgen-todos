package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/maruware/gqlgen-todos/dataloader"
	"github.com/maruware/gqlgen-todos/external"
	"github.com/maruware/gqlgen-todos/graph"
	"github.com/maruware/gqlgen-todos/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := external.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	dataloaderMw := dataloader.DataLoaderMiddleware(db)

	r.Get("/", playground.Handler("GraphQL playground", "/query"))
	r.With(dataloaderMw).Post("/query", srv.ServeHTTP)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
