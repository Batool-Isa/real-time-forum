package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"real-time-forum/backend/database"
	"real-time-forum/backend/handler"
	"real-time-forum/backend/middleware"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"log"
)

func main() {
	file, fileErr := os.OpenFile("LogFile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	log.SetOutput(file)

	err := database.CreateDB("real-forum.db")
	if err != nil {
		log.Fatalf("Error creating or connecting to database: %v", err)
	}

	database.CreateTables()
	database.AddDummyData()
	// Single handler for all routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// Handle API routes
			switch r.URL.Path {
			case "/api/posts":
				handler.GetPostsHandler(w, r)
			case "/api/categories":
				handler.GetCategoriesHandler(w, r)
			case "/api/posts/create":
				middleware.SessionValidator(http.HandlerFunc(handler.CreatePostHandler)).ServeHTTP(w, r)
			case "/api/like":
				middleware.SessionValidator(http.HandlerFunc(handler.LikePost))
			case "/api/dislike":
				middleware.SessionValidator(http.HandlerFunc(handler.DislikePost)).ServeHTTP(w, r)
			case "/api/comment":
				middleware.SessionValidator(http.HandlerFunc(handler.CommentHandler)).ServeHTTP(w, r)
			case "/api/users":
				handler.GetUsersHandler(w, r)
			// case "/api/chat/history/":
			// 	handler.ChatHistoryHandler(w, r)
			default:
				http.NotFound(w, r)
			}
			return
		}
		// Serve index.html for all non-API routes
		http.ServeFile(w, r, "template/index.html")
	})

	//login end point
	http.Handle("/login", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/register", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.RegisterUserHandler)))
	http.Handle("/create", middleware.SessionValidator(http.HandlerFunc(handler.CreateHandler)))
	//http.Handle("/posts", middleware.SessionValidator(http.HandlerFunc(handler.GetPostsHandler)))
	http.Handle("/logout", middleware.SessionValidator(http.HandlerFunc(handler.Logout)))
	http.Handle("/post", middleware.SessionValidator(http.HandlerFunc(handler.PostHandler)))
	http.Handle("/api/saveMessage", middleware.SessionValidator(http.HandlerFunc(handler.SaveMessageHandler)))
	http.Handle("/api/chat/history", middleware.SessionValidator(http.HandlerFunc(handler.GetChatHistoryHandler)))
	http.Handle("/api/user-chat", middleware.SessionValidator(http.HandlerFunc(handler.GetUserChat)))
	//http.Handle("/", middleware.OptionalSessionMiddleware(http.HandlerFunc(handler.SPAHandler)))
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

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.WebSocketHandler(w, r)
	})

	// Static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("template/assets"))))

	fmt.Println("Server started at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
