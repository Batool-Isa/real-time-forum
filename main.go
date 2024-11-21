package main

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	_ "github.com/mattn/go-sqlite3"
	"real-time-forum/backend/handler"
	//"real-time-forum/backend/utils"
	"os"
	"log"
)

func main() {
	file, fileErr := os.OpenFile("myLOG.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	log.SetOutput(file)

	err := database.CreateDB("real-forum.db")
	if err != nil {
		log.Fatalf("Error creating or connecting to database: %v", err)
	}

	// Create tables
	database.CreateTables()
	//go handler.hub.run()// Start the hub in a goroutine

	// Serve static files from the "assets" directory
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("template/assets"))))
    


	// Register handlers for the WebSocket and the main page
	http.HandleFunc("/", handler.IndexHandler)         // Main page handler
	http.HandleFunc("/ws", handler.WebSocketHandler) // WebSocket handler
	http.HandleFunc("/register", handler.RegisterUserHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/create-post", handler.createHandler)


	fmt.Println("Database setup complete")	
	fmt.Println("Server started at http://localhost:8888/")

	// Start the server
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("Error starting server on port 8888:", err)
	}
}
