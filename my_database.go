package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/FilipLusnia/rssagg/internal/database"
	_ "github.com/lib/pq"
)

func myDatabase() *database.Queries {
	dbURL := os.Getenv("DB")
	if dbURL == "" {
		log.Fatal("Database URL is not found")
	}

	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	return database.New(dbConnection)
}
