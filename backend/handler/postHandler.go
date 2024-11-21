package handler

import (
"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/struct"
	//"real-time-forum/backend/utils"
	"database/sql"
	"net/http"
	"strconv"
	"fmt"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	postId := r.URL.Query().Get("id")

	if postId == "" {
		fmt.Println("ERROR")
		///utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(postId)
	if err != nil {
		fmt.Println("ERROR")

//		utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	post, err := database.GetPostById(post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ERROR")

//			utils.ErrorHandler(w, r, http.StatusNotFound)
			return
		}
		fmt.Println("ERROR")
		
		//utils.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := struct {
		Post    structs.Post
		Session *structs.Session
	}{
		Post:    post,
		Session: session,
	}
	fmt.Println("ERROR")

	//utils.RenderTemplate(w, r, "comment.html", data)
}
