package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*structs.Client]bool)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("WebSocket connection attempt...")
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }

    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        log.Println("No session found")
        conn.Close()
        return
    }

    log.Printf("WebSocket connected for user: %d", session.UserID)
    client := &structs.Client{
        Conn:     conn,
        UserID:   session.UserID,
        Username: session.UserName,
    }
    
    clients[client] = true
    go handleMessages(client)
}func handleMessages(client *structs.Client) {
    log.Printf("Client connected: UserID=%d, Username=%s", client.UserID, client.Username)
    defer func() {
        client.Conn.Close()
        delete(clients, client)
    }()

    for {
        _, p, err := client.Conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        var msg structs.Message
        if err := json.Unmarshal(p, &msg); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        // Set sender ID from the authenticated client
        msg.SenderID = client.UserID

        // Add debug logging
        log.Printf("Saving message: Content=%s, SenderID=%d, ReceiverID=%d", 
            msg.Content, msg.SenderID, msg.ReceiverID)

        // Save to database with explicit error handling
        if err := database.SaveMessage(msg.Content, msg.SenderID, msg.ReceiverID); err != nil {
            log.Printf("Failed to save message: %v", err)
        } else {
            log.Printf("Message saved successfully")
        }

        // Forward message to recipient
        for c := range clients {
            if c.UserID == msg.ReceiverID {
                c.Conn.WriteMessage(websocket.TextMessage, p)
                break
            }
        }
    }
}