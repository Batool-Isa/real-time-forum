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
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }

    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        conn.Close()
        return
    }

    client := &structs.Client{
        Conn:     conn,
        UserID:   session.UserID,
        Username: session.UserName,
    }
    
    clients[client] = true

    // Start handling messages in a goroutine
    go handleMessages(client)
}
func handleMessages(client *structs.Client) {
	log.Printf("Client connected ...........")
    defer func() {
        client.Conn.Close()
        delete(clients, client)
    }()

    for {
        messageType, p, err := client.Conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        var msg structs.Message
        if err := json.Unmarshal(p, &msg); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        // Save to database
        err = database.SaveMessage(msg.Content, client.UserID, msg.ReceiverID)
        if err != nil {
            log.Printf("Message saved successfully")
        }

        // Forward message to recipient
        for c := range clients {
            if c.UserID == msg.ReceiverID {
                c.Conn.WriteMessage(messageType, p)
                break
            }
        }
    }
}