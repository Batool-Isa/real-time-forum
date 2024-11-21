package handler

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	//"real-time-forum/backend/utils"
	"real-time-forum/backend/struct"

	"fmt"
	"net/http"
	"strconv"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.SessionKey).(structs.Session)
	if !ok {
		//utils.ErrorHandler(w, r, http.StatusForbidden)
		fmt.Println("error")

		//http.Error(w, "Unable to retrieve session", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		//utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		fmt.Println("error")

		//http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			//utils.ErrorHandler(w, r, http.StatusBadRequest)
			fmt.Println("error")

			return
		}
		postID := r.FormValue("post_id")
		commentText := r.FormValue("comment")

		uid, err := RetrieveLoggedUser(r)
		if err != nil {
			fmt.Println("error")

			//utils.ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		postIDInt, err := strconv.Atoi(postID)
		if err != nil {
		///	utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("error")

			return
		}
		inserterr := database.InsertComment(commentText, uid, postIDInt)
		if inserterr != nil {
		///	utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("error")

		}
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postIDInt), http.StatusSeeOther)
		return
	}
	fmt.Println("error")
	//utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)
}
