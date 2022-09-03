package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fluffy-octo/ik-reddit-backend/engine/admins"
	"github.com/fluffy-octo/ik-reddit-backend/graph"
	"github.com/fluffy-octo/ik-reddit-backend/graph/generated"
	"github.com/fluffy-octo/ik-reddit-backend/utils"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8081"

func main() {
	utils.InitialiseDB()
	admins.AddAdmin()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()
	// router.Use(auth.ValidateBasicAuth())
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler)

	graphQlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", graphQlServer)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Println("🚀 Server ready at http://localhost:" + port + "/query")

	log.Fatal(srv.ListenAndServe())
}
