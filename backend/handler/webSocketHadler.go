package handler

import (
	//"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader with custom configuration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { 
		 return true }, // Allow all origins
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    defer conn.Close()

    // Start a ping loop to keep the connection alive
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        for {
            <-ticker.C
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Println("WebSocket ping error:", err)
                return
            }
        }
    }()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("WebSocket read error:", err)
            break
        }
        log.Printf("Received: %s", message)
        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Println("WebSocket write error:", err)
            break
        }
    }
}

