package handler

import (
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/structs"
	"sync"
	"fmt"
	"encoding/json"
	"time"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*structs.Client]bool)
	broadcast = make(chan structs.Message)
	mutex     = &sync.Mutex{}
	activeUsers = make(map[int]*structs.Client)
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
		LastActivity: time.Now(),
	}

	mutex.Lock()
	clients[client] = true
	activeUsers[userID] = client
	fmt.Println(activeUsers[userID])
	mutex.Unlock()
	// for userID, _ := range activeUsers {
	// 	// Notify others that this user is online
	// 	notifyPresenceChange(userID, true) 
	// }
	

	defer func() {
		mutex.Lock()
		delete(clients, client)
		//remove user froma ctive users map
		delete(activeUsers, client.UserID)
		mutex.Unlock()
		//Notify the users that this user is offline
		// notifyPresenceChange(client.UserID, false)
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

// func notifyPresenceChange(userID int, online bool) {
//     presence := structs.Message{
//         MessageID:  0,
//         Content:    fmt.Sprintf("User %d is now %s", userID, onlineStatus(online)),
//         SenderID:   userID,
//         ReceiverID: 0, // Broadcast to all users
//     }

//     mutex.Lock()
//     for _, client := range activeUsers {
//         err := client.Conn.WriteJSON(presence)
//         if err != nil {
//             log.Printf("Error sending presence update: %v", err)
//             client.Conn.Close()
//             delete(activeUsers, client.UserID) // Remove disconnected client
//         }
//     }
//     mutex.Unlock()
// }


func onlineStatus(online bool) string {
    if online {
        return "online"
    }
    return "offline"
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

func GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
    currentUserID, err := RetrieveLoggedUser(r)
    if err != nil {
        log.Printf("Error retrieving current user: %v", err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    mutex.Lock()
    defer mutex.Unlock()

    var onlineUsers []structs.User
    for userID := range activeUsers {
        if userID == currentUserID {
            continue // Skip the current user
        }
        user, err := database.GetUserByID(userID)
        if err != nil {
            log.Printf("Error fetching user details for userID %d: %v", userID, err)
            continue
        }
        onlineUsers = append(onlineUsers, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(onlineUsers)
}

