package main

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/handler"
	"real-time-forum/backend/middleware"

	_ "github.com/mattn/go-sqlite3"

	"log"
)

func main() {
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
	http.HandleFunc("/", handler.IndexHandler)       // Main page handler
	http.HandleFunc("/ws", handler.WebSocketHandler) // WebSocket handler
	http.HandleFunc("/register", handler.RegisterUserHandler)
	// http.Handle("/login", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.LoginHandler)))
	// http.Handle("/register", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.RegisterUserHandler)))
	http.Handle("/create_post", middleware.SessionValidator(http.HandlerFunc(handler.CreateHandler)))
	http.Handle("/like", middleware.SessionValidator(http.HandlerFunc(handler.LikePost)))
	http.Handle("/dislike", middleware.SessionValidator(http.HandlerFunc(handler.DislikePost)))
	http.Handle("/logout", middleware.SessionValidator(http.HandlerFunc(handler.Logout)))
	http.Handle("/post", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.PostHandler)))
	http.Handle("/add_comment", middleware.SessionValidator(http.HandlerFunc(handler.CommentHandler)))
	http.Handle("/like_comment", middleware.SessionValidator(http.HandlerFunc(handler.LikeComment)))
	http.Handle("/dislike_comment", middleware.SessionValidator(http.HandlerFunc(handler.DislikeComment)))
	http.HandleFunc("/api/users", handler.GetUsersHandler)
	http.Handle("/api/chat/history", middleware.SessionValidator(http.HandlerFunc(handler.GetChatHistoryHandler)))


	fmt.Println("Database setup complete")
	fmt.Println("Server started at http://localhost:8080/")

	// Start the server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server on port 8080:", err)
	}
}
