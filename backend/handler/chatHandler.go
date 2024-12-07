package handler

import (
	"encoding/json"
	"fmt"

	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"

	//"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"
	"strconv"
)

func GetUserChat(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userID := session.UserID
    users, err := database.GetAllChats(userID)
    if err != nil {
        http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
        return
    }
	// In case the user has not chatted with anyone yet, return an empty array
	if users == nil {
		users = []structs.User{}
	}
fmt.Println("Debug", users)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}


func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get session from context
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	senderID := session.UserID // Get current user's ID from session
	receiverID := r.URL.Query().Get("receiverId")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit := 10 

	if receiverID == "" {
		http.Error(w, "ReceiverID is required", http.StatusBadRequest)
		return
	}

	receiverIDInt, _ := strconv.Atoi(receiverID)

	// Fetch chat history
	messages, err := database.GetChatHistory(senderID, receiverIDInt, offset, limit)
	if err != nil {
		http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
		return
	}

	if messages == nil {
        messages = []structs.Message{} // Ensure empty array, not null
    }


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func SaveMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var msg structs.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	msg.SenderID = session.UserID

	// Save to database
	err = database.SaveMessage(msg.Content, msg.SenderID, msg.ReceiverID)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Broadcast message through WebSocket
	broadcast <- msg

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message saved successfully"})
}