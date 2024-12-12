package handler

import (
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/structs"
	"sync"
	//"strconv"
	//"strings"
	"fmt"
	"encoding/json"
	"time"
	"github.com/gorilla/websocket"
	//	"github.com/gorilla/mux"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}
type WebSocketMessage struct {
	Type    string      `json:"type"`    // The type of message (e.g., "chat", "presence")
	Payload interface{} `json:"payload"` // The actual message payload
}

var (
	clients   = make(map[*structs.Client]bool)
	 broadcast = make(chan WebSocketMessage)
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
	
	// Notify others that this user is online
	go notifyPresenceChange(userID, true)

	defer func() {
		mutex.Lock()
		delete(clients, client)
		//remove user froma ctive users map
		delete(activeUsers, client.UserID)
		mutex.Unlock()
	
		// Notify others that this user is offline
		notifyPresenceChange(userID, false)

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
		broadcast <- WebSocketMessage{
			Type:    "chat",
			Payload: msg,
		}
	}
	
}

func notifyPresenceChange(userID int, online bool) {
    message := WebSocketMessage{
        Type: "presence",
        Payload: struct {
            UserID int  `json:"userId"`
            Online bool `json:"online"`
        }{
            UserID: userID,
            Online: online,
        },
    }

    mutex.Lock()
    for client := range clients {
        err := client.Conn.WriteJSON(message)
        if err != nil {
            log.Printf("Error notifying presence change: %v", err)
            client.Conn.Close()
            delete(clients, client)
        }
    }
    mutex.Unlock()
}


func onlineStatus(online bool) string {
    if online {
        return "online"
    }
    return "offline"
}


func handleMessages() {
	for msg := range broadcast {
		// Save the message to the database if it's a chat message
		if msg.Type == "chat" {
			err := database.SaveMessage(msg.Payload.(structs.Message).Content, msg.Payload.(structs.Message).SenderID, msg.Payload.(structs.Message).ReceiverID)
			if err != nil {
				log.Printf("Error saving message: %v", err)
				continue
			}
		}

		// Broadcast the message to relevant clients
		mutex.Lock()
		for client := range clients {
			if msg.Type == "chat" && (client.UserID == msg.Payload.(structs.Message).ReceiverID || client.UserID == msg.Payload.(structs.Message).SenderID) {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					client.Conn.Close()
					delete(clients, client)
				}
			} else if msg.Type == "presence" {
				// Broadcast presence updates to all connected clients
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending presence update: %v", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}
}


func GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    defer mutex.Unlock()

    var onlineUsers []int
    for userID := range activeUsers {
        onlineUsers = append(onlineUsers, userID)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(onlineUsers)
}

func GetAllUsersWithStatus(w http.ResponseWriter, r *http.Request) {
    currentUserID, err := RetrieveLoggedUser(r)
    if err != nil {
        log.Printf("Error retrieving current user: %v", err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Fetch all users from the database
    users, err := database.FetchAllUsers()
    if err != nil {
        log.Printf("Error fetching all users: %v", err)
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    mutex.Lock()
    defer mutex.Unlock()

    // Create a slice for users with their online status
    var usersWithStatus []struct {
        UserID   int    `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
        Username string `json:"username"`
        Online   bool   `json:"online"`
    }

    for _, user := range users {
        // Skip the current user
        if user.UserID == currentUserID {
            continue
        }

        // Check if the user is online
        _, isOnline := activeUsers[user.UserID]

        usersWithStatus = append(usersWithStatus, struct {
            UserID   int    `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
            Username string `json:"username"`
            Online   bool   `json:"online"`
        }{
            UserID:   user.UserID,
			FirstName: user.FirstName,
			LastName: user.LastName,
            Username: user.Username,
            Online:   isOnline,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(usersWithStatus)
}



// func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
//     // Extract userId from the URL
//     pathParts := strings.Split(r.URL.Path, "/")
//     if len(pathParts) < 4 { // Assuming URL is like /api/user/{userId}
//         http.Error(w, "Invalid URL", http.StatusBadRequest)
//         return
//     }

//     userID := pathParts[3] // Get the ID from the URL
// 	userIDInt, err := strconv.Atoi(userID)
//     // Fetch the user from the database
//     user, err := database.GetUserByID(userIDInt)
//     if err != nil {
//         http.Error(w, "User not found", http.StatusNotFound)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(user)
// }
