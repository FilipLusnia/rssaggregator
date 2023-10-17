package main

import (
	"log"
	"net/http"
	"os"

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

	apiCfg := apiConfig{
		DB: myDatabase(),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/status", handlerReadiness)
	router.Get("/error", handlerError)
	router.Post("/user", apiCfg.handlerCreateUser)

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
