package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var postData struct {
        UserID      int      `json:"user_id"`
        PostHeading string   `json:"post_heading"`
        PostData    string   `json:"post_data"`
        Categories  []string `json:"categories"`
    }

    if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := database.InsertPost(postData.UserID, postData.PostHeading, postData.PostData, postData.Categories)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var likeData struct {
        PostID int `json:"post_id"`
        UserID int `json:"user_id"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&likeData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := database.InsertLikes(likeData.PostID, likeData.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var dislikeData struct {
        PostID int `json:"post_id"`
        UserID int `json:"user_id"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&dislikeData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := database.InsertDislikes(dislikeData.PostID, dislikeData.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
func CommentHandlerPost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var commentData struct {
        Comment string `json:"comment"`
        UserID  int    `json:"user_id"`
        PostID  int    `json:"post_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&commentData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := database.InsertComment(commentData.Comment, commentData.UserID, commentData.PostID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    users, err := database.FetchAllUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)

}
func ChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Extract user ID from URL path
    userIDStr := r.URL.Path[len("/api/chat/history/"):]
    
    // Convert userID from string to int
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Assuming the second argument is the current user's ID from the session
    currentUserID := session.UserID

    messages, err := database.GetChatHistory(userID, currentUserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}