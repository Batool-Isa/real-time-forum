package main

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
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
	database.AddDummyData()
	//go handler.hub.run()// Start the hub in a goroutine

	// Serve static files from the "assets" directory
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("template/assets"))))
    


	// Register handlers for the WebSocket and the main page
	// http.HandleFunc("/", handler.IndexHandler)         // Main page handler
	// http.HandleFunc("/ws", handler.WebSocketHandler) // WebSocket handler
	// http.HandleFunc("/register", handler.RegisterUserHandler)
	// http.HandleFunc("/login", handler.LoginHandler)
	// http.HandleFunc("/create-post", handler.CreateHandler)

	http.Handle("/", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.IndexHandler)))
	http.Handle("/login", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/register", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.RegisterUserHandler)))
	http.Handle("/create-post", middleware.SessionValidator(http.HandlerFunc(handler.CreateHandler)))
	http.Handle("/like", middleware.SessionValidator(http.HandlerFunc(handler.LikePost)))
	http.Handle("/ws", middleware.SessionValidator(http.HandlerFunc(handler.WebSocketHandler)))
	http.Handle("/dislike", middleware.SessionValidator(http.HandlerFunc(handler.DislikePost)))
	http.Handle("/logout", middleware.SessionValidator(http.HandlerFunc(handler.Logout)))
	http.Handle("/posts", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.GetPostsHandler)))
	http.Handle("/add_comment", middleware.SessionValidator(http.HandlerFunc(handler.CommentHandler)))
	http.Handle("/like_comment", middleware.SessionValidator(http.HandlerFunc(handler.LikeComment)))
	http.Handle("/dislike_comment", middleware.SessionValidator(http.HandlerFunc(handler.DislikeComment)))

	fmt.Println("Database setup complete")	
	fmt.Println("Server started at http://localhost:8888/")

	// Start the server
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("Error starting server on port 8888:", err)
	}
}
