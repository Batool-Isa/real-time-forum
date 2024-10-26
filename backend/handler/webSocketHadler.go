package handler
// import (
//     "log"
//     "net/http"

//     "github.com/gorilla/websocket"
// )

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

//     // Create client instance
//     client := &Client{
//         conn: conn,
//         send: make(chan []byte),
//     }

//     // Register client
//     hub.register <- client

//     // Start goroutines for reading and writing messages
//     go client.writePump()
//     go client.readPump()
// }