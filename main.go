package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"real-time-forum/backend/database"
	"real-time-forum/backend/handler"
	"real-time-forum/backend/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Configure logging
	file, fileErr := os.OpenFile("LogFile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	log.SetOutput(file)

	// Initialize database
	err := database.CreateDB("real-forum.db")
	if err != nil {
		log.Fatalf("Error creating or connecting to database: %v", err)
	}

	database.CreateTables()

	// Register routes
	registerRoutes()

	// Start the server
	fmt.Println("Server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerRoutes() {
	// Static files (e.g., CSS, JS, Images)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("template/assets"))))

	// WebSocket endpoint
	http.HandleFunc("/ws", handler.WebSocketHandler)

	// Authentication routes
	http.Handle("/login", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/register", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.RegisterUserHandler)))
	http.Handle("/logout", middleware.SessionValidator(http.HandlerFunc(handler.Logout)))

	// Post-related routes
	http.Handle("/create", middleware.SessionValidator(http.HandlerFunc(handler.CreateHandler)))
	http.Handle("/api/posts", middleware.SessionValidator(http.HandlerFunc(handler.IndexHandler)))

	// API routes
	//http.HandleFunc("/api/posts", handler.GetPostsHandler)
	http.HandleFunc("/api/categories", handler.GetCategoriesHandler)
	http.Handle("/api/posts/create", middleware.SessionValidator(http.HandlerFunc(handler.CreatePostHandler)))
	http.Handle("/api/posts/like", middleware.SessionValidator(http.HandlerFunc(handler.LikePost)))
	http.Handle("/api/posts/dislike", middleware.SessionValidator(http.HandlerFunc(handler.DislikePost)))
	http.Handle("/api/comments", middleware.SessionValidator(http.HandlerFunc(handler.CommentHandler)))
	http.HandleFunc("/api/users", handler.GetUsersHandler)
	http.HandleFunc("/api/chat/history/", handler.ChatHistoryHandler)
	http.HandleFunc("/api/check-session", CheckSessionHandler)


	// Root route to decide whether to show login or redirect to posts
	http.HandleFunc("/", rootHandler)
}

// Root handler: Redirect to login or posts based on session
func rootHandler(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		// If no session, serve login page
		http.ServeFile(w, r, "template/index.html")
		return
	}
	// If session exists, redirect to posts
	http.Redirect(w, r, "/api/posts", http.StatusSeeOther)
}
func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
    session := middleware.GetSessionFromContext(r.Context())
    if session != nil {
        json.NewEncoder(w).Encode(map[string]bool{"isLoggedIn": true})
    } else {
        json.NewEncoder(w).Encode(map[string]bool{"isLoggedIn": false})
    }
}

