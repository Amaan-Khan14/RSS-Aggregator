package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// feed, err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

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

	go startScrapping(apiCfg.DB, 10, time.Minute) //go routine

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
	v1Router.Get("/get", apiCfg.middlewareAuht(apiCfg.handlerGetUser))

	//get user by name route
	//path: "/v1/get/name"
	v1Router.Get("/get/name", apiCfg.handlerGetUserByName)

	//create user route
	//path: "/v1/create"
	v1Router.Post("/create", apiCfg.handlerCreateUsers)

	//creaet feed route
	//path: "/v1/create/feed"
	v1Router.Post("/create/feed", apiCfg.middlewareAuht(apiCfg.handlerCreateFeed))

	//get feed route
	//path: "/v1/get/feed"
	v1Router.Get("/get/feed", apiCfg.handlerGetFeeds)

	//feedfollow route
	//path: "/v1/create/feedfollow"
	v1Router.Post("/create/feedfollow", apiCfg.middlewareAuht(apiCfg.handlerCreateFeedFollow))

	//get feedfollow route
	//path: "/v1/get/feedfollow"
	v1Router.Get("/get/feedfollow", apiCfg.middlewareAuht(apiCfg.getFeedFollows))

	//delete feedfollow route
	//path: "/v1/delete/feedfollow"
	v1Router.Delete("/delete/feedfollow/{feedFollowId}", apiCfg.middlewareAuht(apiCfg.deleteFeedFollow))

	//get posts route
	//path: "/v1/get/posts"
	v1Router.Get("/get/posts", apiCfg.middlewareAuht(apiCfg.handlerGetPostsForUser))

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
