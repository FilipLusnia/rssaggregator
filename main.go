package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/FilipLusnia/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	dbConn := myDatabase()
	apiCfg := apiConfig{
		DB: dbConn,
	}

	go startScraping(dbConn, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ---------------------

	router.Get("/status", handlerReadiness)
	router.Get("/error", handlerError)

	router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	router.Post("/users", apiCfg.handlerCreateUser)

	router.Get("/feeds", apiCfg.handlerGetFeeds)
	router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Get("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	router.Delete("/feed-follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	// ---------------------

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("Server starting on port %v", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
