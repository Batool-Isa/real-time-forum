package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"
)



func LikePost(w http.ResponseWriter, r *http.Request) {
    session, ok := r.Context().Value(middleware.SessionKey).(structs.Session)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
        return
    }

    var requestData struct {
        PostID int `json:"post_id"`
    }

    // Read and log the request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Println("Error reading body:", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    log.Println("Raw request body:", string(body))

    // Unmarshal the JSON data into requestData
    err = json.Unmarshal(body, &requestData)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
        return
    }

    log.Println("Parsed PostID:", requestData.PostID)

    // Check if PostID is valid
    if requestData.PostID == 0 {
        log.Println("PostID is missing or invalid")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Missing or invalid post_id"})
        return
    }

    uid := session.UserID
    log.Printf("UserID: %d is liking PostID: %d", uid, requestData.PostID)

    // Perform the like operation
    err = database.InsertLikes(requestData.PostID, uid)
    if err != nil {
        log.Println("Error inserting like:", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to like post"})
        return
    }
    database.DeleteDislike(requestData.PostID, uid)
    // Update the post counts
    post, err := database.UpdatePost(requestData.PostID)
    if err != nil {
        log.Println("Error updating post:", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update post counts"})
        return
    }

    // Send the updated post data as JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(post)
}


func DislikePost(w http.ResponseWriter, r *http.Request) {
    session, ok := r.Context().Value(middleware.SessionKey).(structs.Session)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
        return
    }

    var requestData struct {
        PostID int `json:"post_id"`
    }

    // Read and log the request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Println("Error reading body:", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    log.Println("Raw request body:", string(body)) // This will help debug if the body is coming through correctly

    // Unmarshal the JSON data into requestData
    err = json.Unmarshal(body, &requestData)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
        return
    }

    log.Println("Parsed PostID:", requestData.PostID)  // Log the parsed PostID

    // Check if PostID is valid
    if requestData.PostID == 0 {
        log.Println("PostID is missing or invalid")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Missing or invalid post_id"})
        return
    }

    uid := session.UserID
    log.Printf("UserID: %d is disliking PostID: %d", uid, requestData.PostID)

    // Perform the dislike operation
    err = database.InsertDislikes(requestData.PostID, uid)
    if err != nil {
        log.Println("Error inserting dislike:", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to dislike post"})
        return
    }
    database.DeleteLike(requestData.PostID, uid)

    // Update the post counts
    post, err := database.UpdatePost(requestData.PostID)
    if err != nil {
        log.Println("Error updating post:", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update post counts"})
        return
    }

    // Send the updated post data as JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(post)
}
