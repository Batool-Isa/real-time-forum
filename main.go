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
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Log file setup
	file, fileErr := os.OpenFile("LogFile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	log.SetOutput(file)

	// Database initialization
	err := database.CreateDB("real-forum.db")
	if err != nil {
		log.Fatalf("Error creating or connecting to database: %v", err)
	}

	database.CreateTables()
	database.AddDummyData()

	// API Endpoints
	http.HandleFunc("/api/posts", handler.GetPostsHandler)
	http.HandleFunc("/api/categories", handler.GetCategoriesHandler)
	http.Handle("/api/posts/create", middleware.SessionValidator(http.HandlerFunc(handler.CreatePostHandler)))
	http.Handle("/api/like", middleware.SessionValidator(http.HandlerFunc(handler.LikePost)))
	http.Handle("/api/dislike", middleware.SessionValidator(http.HandlerFunc(handler.DislikePost)))
	http.Handle("/api/comment", middleware.SessionValidator(http.HandlerFunc(handler.CommentHandler)))
	http.HandleFunc("/api/users", handler.GetUsersHandler)
	http.Handle("/api/saveMessage", middleware.SessionValidator(http.HandlerFunc(handler.SaveMessageHandler)))
	http.Handle("/api/chat/history", middleware.SessionValidator(http.HandlerFunc(handler.GetChatHistoryHandler)))
	http.Handle("/api/user-chat", middleware.SessionValidator(http.HandlerFunc(handler.GetUserChat)))

	http.HandleFunc("/api/session-status", func(w http.ResponseWriter, r *http.Request) {
		userID, err := handler.RetrieveLoggedUser(r)
		if err != nil {
			http.Error(w, "No active session", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"userId": userID,
		})
		w.WriteHeader(http.StatusOK)
	})

	// Authentication Endpoints
	http.Handle("/login", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/register", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.RegisterUserHandler)))
	http.Handle("/logout", middleware.SessionValidator(http.HandlerFunc(handler.Logout)))

	// Post Management Endpoints
	http.Handle("/create", middleware.SessionValidator(http.HandlerFunc(handler.CreateHandler)))
	http.Handle("/post", middleware.SessionValidator(http.HandlerFunc(handler.PostHandler)))

	// WebSocket Endpoint
	http.HandleFunc("/ws", handler.WebSocketHandler)

	// Static Files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("template/assets"))))

	// Fallback route for serving the SPA
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "template/index.html")
	})

	fmt.Println("Server started at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
