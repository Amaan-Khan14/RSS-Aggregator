package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("port not found")
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handlerReadiness)
	//Full path will be "/v1/readiness"

	//error route
	v1Router.Get("/error", hanlErr)
	//Full path will be "/v1/error"

	router.Mount("/v1", v1Router)

	log.Print("Server is running on port: ", port)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}
