package handler

import (
	"encoding/json"
	
	"net/http"
	"real-time-forum/backend/database"
)


func CommentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
        return
    }

    // Retrieve the logged-in user's ID
    uid, err := RetrieveLoggedUser(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    // Retrieve the username for the logged-in user
    username, err := database.GetUsername(uid)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve username"})
        return
    }

    // Parse the JSON body
    var reqData struct {
        PostID      int    `json:"post_id"`
        CommentText string `json:"comment"`
    }

    err = json.NewDecoder(r.Body).Decode(&reqData)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    // Validate input
    if reqData.PostID == 0 || reqData.CommentText == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Missing post ID or comment text"})
        return
    }

    // Insert the comment into the database
    err = database.InsertComment(reqData.CommentText, uid, reqData.PostID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add comment"})
        return
    }

    // Respond with success and include the username
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message":  "Comment added successfully",
        "username": username,
    })
}
