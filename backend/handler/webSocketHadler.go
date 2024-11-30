package handler

import (
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/structs"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*structs.Client]bool)
	broadcast = make(chan structs.Message)
	mutex     = &sync.Mutex{}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	userID, err := RetrieveLoggedUser(r)
	if err != nil {
		log.Println("Error getting user ID:", err)
		conn.Close()
		return
	}

	client := &structs.Client{
		Conn:   conn,
		UserID: userID,
	}

	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		conn.Close()
	}()

	// Start goroutine for handling messages
	go handleMessages()

	// Read messages from client
	for {
		var msg structs.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		msg.SenderID = userID
		broadcast <- msg
	}
}

func handleMessages() {
	for msg := range broadcast {
		// Save to database first
		err := database.SaveMessage(msg.Content, msg.SenderID, msg.ReceiverID)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}
	
		// Then broadcast to recipient
		mutex.Lock()
		for client := range clients {
			if client.UserID == msg.ReceiverID {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}
}
