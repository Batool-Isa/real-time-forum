package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/structs"
    "fmt"
    "time"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*structs.Client]bool)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Send and receive messages with the client
	for {
		// Read message from WebSocket
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		fmt.Printf("Received message: %s\n", message)

		// Example response: send the current time
		response := fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC3339))
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}
