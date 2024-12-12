package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

)



func LikePostTest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var requestData struct {
        PostID int `json:"post_id"`
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Println("Error reading body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Println("Raw request body:", string(body))

    err = json.Unmarshal(body, &requestData)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    log.Printf("Parsed PostID: %d", requestData.PostID)

    // Return a static response for testing
    response := map[string]interface{}{
        "status":  "success",
        "post_Id": requestData.PostID,
        "likes":   10, // Mocked data
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
