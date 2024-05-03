package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("port not found")
	}

	//database connection
	dbURL := os.Getenv("DB_URI")
	if dbURL == "" {
		log.Fatal("database url not found")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Cant communicate with database", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
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

	//path: "/v1/readiness"
	v1Router.Get("/readiness", handlerReadiness)

	//error route
	//path: "/v1/error"
	v1Router.Get("/error", hanlErr)

	//get user by apikey route
	//path: "/v1/get"
	v1Router.Get("/get", apiCfg.handlerGetUser)

	//get user by name route
	//path: "/v1/get/name"
	v1Router.Get("/get/name", apiCfg.handlerGetUserByName)

	//create user route
	//path: "/v1/create"
	v1Router.Post("/create", apiCfg.handlerCreateUsers)

	router.Mount("/v1", v1Router)

	log.Print("Server is running on port: ", port)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}

type apiConfig struct {
	DB *database.Queries
}
