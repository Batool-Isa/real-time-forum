package handler

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/struct"
	//"real-time-forum/backend/utils"
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category")

	var posts []structs.Post
	var err error

	if categoryID == "" || categoryID == "all" {
		posts, err = database.GetAllPosts()
	} else {
		catID, err := strconv.Atoi(categoryID)
		if err != nil {
			//utils.ErrorHandler(w, r, http.StatusBadRequest)
			fmt.Println("error")

			//http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		posts, err = database.GetPostsByCategory(catID)
		if err != nil {
			fmt.Println("error")
			
			//utils.ErrorHandler(w, r, http.StatusInternalServerError)
			//http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}
	}

	if err != nil {
		fmt.Println("error")
		
		//utils.ErrorHandler(w, r, http.StatusInternalServerError)
		//http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
