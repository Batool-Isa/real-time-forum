// package main

// import (
//     "log"
//     "net/http"
//     "sync"

//     "github.com/gorilla/websocket"
// )

// // Client represents a single WebSocket connection.
// type Client struct {
//     conn *websocket.Conn // The WebSocket connection
//     send chan []byte     // Channel for sending messages to the client
// }

// // Hub maintains the set of active clients and broadcasts messages to them.
// type Hub struct {
//     clients    map[*Client]bool         // Registered clients
//     register   chan *Client             // Channel for registering new clients
//     unregister chan *Client             // Channel for unregistering clients
//     broadcast  chan []byte              // Channel for broadcasting messages
//     mu         sync.Mutex               // Mutex for concurrent access
// }

// // Initialize a new hub instance
// var hub = Hub{
//     clients:    make(map[*Client]bool),
//     register:   make(chan *Client),
//     unregister: make(chan *Client),
//     broadcast:  make(chan []byte),
// }

// // Run the hub to manage client connections and broadcasts
// func (h *Hub) run() {
//     for {
//         select {
//         case client := <-h.register:
//             h.mu.Lock()
//             h.clients[client] = true // Add the new client
//             h.mu.Unlock()

//         case client := <-h.unregister:
//             h.mu.Lock()
//             if _, ok := h.clients[client]; ok {
//                 delete(h.clients, client) // Remove the client
//                 close(client.send)        // Close the send channel
//             }
//             h.mu.Unlock()

//         case message := <-h.broadcast:
//             h.mu.Lock()
//             for client := range h.clients {
//                 select {
//                 case client.send <- message: // Send the message to the client
//                 default:
//                     close(client.send) // If the client is not ready, close the channel
//                     delete(h.clients, client) // Remove the client
//                 }
//             }
//             h.mu.Unlock()
//         }
//     }
// }

// // WebSocket handler to manage incoming WebSocket connections
// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
//     upgrader := websocket.Upgrader{
//         ReadBufferSize:  1024,
//         WriteBufferSize: 1024,
//     }
//     conn, err := upgrader.Upgrade(w, r, nil)
//     if err != nil {
//         log.Printf("WebSocket upgrade failed: %v", err)
//         return
//     }
//     defer conn.Close()

//     // Create a new client instance
//     client := &Client{
//         conn: conn,
//         send: make(chan []byte),
//     }

//     // Register the client in the hub
//     hub.register <- client

//     // Start goroutines for reading and writing messages
//     go client.writePump()
//     go client.readPump()
// }

// // Read messages from the WebSocket connection
// func (c *Client) readPump() {
//     defer func() {
//         hub.unregister <- c
//         c.conn.Close()
//     }()
//     for {
//         _, message, err := c.conn.ReadMessage()
//         if err != nil {
//             log.Printf("Error reading message: %v", err)
//             break
//         }
//         // Broadcast the received message to all clients
//         hub.broadcast <- message
//     }
// }

// // Write messages to the WebSocket connection
// func (c *Client) writePump() {
//     defer c.conn.Close()
//     for message := range c.send {
//         err := c.conn.WriteMessage(websocket.TextMessage, message)
//         if err != nil {
//             log.Printf("Error writing message: %v", err)
//             break
//         }
//     }
// }

// // Main function to start the server
// func main() {
//     go hub.run() // Start the hub in a goroutine

//     http.HandleFunc("/ws", handleWebSocket) // Handle WebSocket connections

//     log.Println("Server is running on :8080")
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         log.Fatalf("Server failed: %v", err)
//     }
// }
