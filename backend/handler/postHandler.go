package handler

import (
	"encoding/json"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	//"real-time-forum/backend/struct"

	//"real-time-forum/backend/utils"
	"database/sql"
	//"fmt"
	"net/http"
	"strconv"
)
func PostHandler(w http.ResponseWriter, r *http.Request) {
    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    postId := r.URL.Query().Get("id")
    if postId == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    post_id, err := strconv.Atoi(postId)
    if err != nil {
        http.Error(w, "Invalid Post ID", http.StatusBadRequest)
        return
    }

    post, err := database.GetPostById(post_id)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Post not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
        return
    }

    // Send the post as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}



// func PostHandler(w http.ResponseWriter, r *http.Request) {
// 	session := middleware.GetSessionFromContext(r.Context())
// 	postId := r.URL.Query().Get("id")

// 	if postId == "" {
// 		fmt.Println("ERROR")
// 		///utils.ErrorHandler(w, r, http.StatusBadRequest)
// 		return
// 	}

// 	post_id, err := strconv.Atoi(postId)
// 	if err != nil {
// 		fmt.Println("ERROR")

// //		utils.ErrorHandler(w, r, http.StatusBadRequest)
// 		return
// 	}

// 	post, err := database.GetPostById(post_id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("ERROR")

// //			utils.ErrorHandler(w, r, http.StatusNotFound)
// 			return
// 		}
// 		fmt.Println("ERROR")
		
// 		//utils.ErrorHandler(w, r, http.StatusInternalServerError)
// 		return
// 	}

// 	data := struct {
// 		Post    structs.Post
// 		Session *structs.Session
// 	}{
// 		Post:    post,
// 		Session: session,
// 	}
// 	fmt.Println("ERROR")

// 	//utils.RenderTemplate(w, r, "comment.html", data)
// }
