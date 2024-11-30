package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"

	//"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"
	"strconv"
)

func GetAllChatsHandler(w http.ResponseWriter, r *http.Request) {
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

    userID := session.UserID // Get current user's ID from session

    // Fetch all chats for the user
    chats, err := database.GetAllChats(userID) // Implement this function in your database package
    if err != nil {
        http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chats)
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

    if receiverID == "" {
        http.Error(w, "ReceiverID is required", http.StatusBadRequest)
        return
    }

    receiverIDInt, _ := strconv.Atoi(receiverID)

    // Fetch chat history
    messages, err := database.GetChatHistory(senderID, receiverIDInt)
    if err != nil {
        http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

func SaveMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Get session from context
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

    // Set SenderID from session
    msg.SenderID = session.UserID

    // Save the message
    err = database.SaveMessage(msg.Content, msg.SenderID, msg.ReceiverID)
    if err != nil {
        http.Error(w, "Failed to save message", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Message saved successfully"})
}
