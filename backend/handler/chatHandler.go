package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"strconv"
)

func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from session
	currentUser := middleware.GetSessionFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get other user ID from query parameter
	otherUserID := r.URL.Query().Get("userId")
	otherID, err := strconv.Atoi(otherUserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch chat history from database
	messages, err := database.GetChatHistory(currentUser.UserID, otherID)
	if err != nil {
		http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
