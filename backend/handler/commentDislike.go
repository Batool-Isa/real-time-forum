
package handler

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/struct"
	//"real-time-forum/backend/utils"
	"fmt"
	"net/http"
	"strconv"
	// "structs"
)

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.SessionKey).(structs.Session)
	if !ok {
		//utils.ErrorHandler(w, r, http.StatusInternalServerError)
		fmt.Println("error")
		return
	}
	if r.Method != http.MethodPost {
		//utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		fmt.Println("error")

		return
	}

	commentId := r.FormValue("comment_id")
	if commentId == "" {
	
		fmt.Println("error")
		// utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post_id")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		//utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("error")

		return
	}

	cid, err := strconv.Atoi(commentId)
	if err != nil {
		fmt.Println("error")
		//utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	uid, err := RetrieveLoggedUser(r)
	if err != nil {
		fmt.Println("error")

		//utils.ErrorHandler(w, r, http.StatusInternalServerError)

		return
	}

	// Insert into dislikes table
	inserterr := database.InsertCommentDislikes(cid, uid)
	if inserterr != nil {
		//utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("error")

	}
	
	database.DeleteCommentLike(cid, uid)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postIDInt), http.StatusSeeOther)
}
