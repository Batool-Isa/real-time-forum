package handler

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	//"real-time-forum/backend/structs"
	//"real-time-forum/backend/utils"
	"strconv"
)


func PostHandler(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		fmt.Println(session)
		fmt.Println(("Session not found"))
	}
    postIdStr := r.URL.Query().Get("id")
    if postIdStr == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }
	postId , err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Post ID can't convert from string to int", http.StatusNotFound)
        return
	}
    post, err := database.GetPostById(postId)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}
