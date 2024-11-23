package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"strconv"
)

func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]structs.Message{})
		return
	}
	
	otherUserID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]structs.Message{})
		return
	}

	currentUser := middleware.GetSessionFromContext(r.Context())
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]structs.Message{})
		return
	}

	messages, err := database.GetChatHistory(currentUser.UserID, otherUserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]structs.Message{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}