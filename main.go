package main

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	_ "github.com/mattn/go-sqlite3"
	"real-time-forum/backend/handler"
	"log"
)

func main() {
	err := database.CreateDB("real-forum.db")
	if err != nil {
		log.Fatalf("Error creating or connecting to database: %v", err)
	}

	// Create tables
	database.CreateTables()

	// Register the IndexHandler function
	http.HandleFunc("/", handler.IndexHandler)

	fmt.Println("Database setup complete")	
	fmt.Println("Server started at http://localhost:8888/")

	// Start the server
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("Error starting server on port 8888:", err)
	}
}
