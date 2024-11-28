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
	
	log.Println("Parsed PostID:", requestData.PostID)
	
    err = json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil || requestData.PostID == 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

	
    uid := session.UserID

    err = database.InsertDislikes(requestData.PostID, uid)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to dislike post"})
        return
    }

    _ = database.DeleteDislike(requestData.PostID, uid)
    post, err := database.UpdatePost(requestData.PostID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update post counts"})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(post)
}

