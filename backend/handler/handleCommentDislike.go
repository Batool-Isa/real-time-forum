package handler

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"
	"real-time-forum/backend/utils"
	"strconv"
	// "structs"
)

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.SessionKey).(structs.Session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		w.WriteHeader(http.StatusForbidden)
		return
		// utils.ErrorHandler(w, r, http.StatusInternalServerError)
		// return
	}
	if r.Method != http.MethodPost {
		utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	commentId := r.FormValue("comment_id")
	if commentId == "" {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post_id")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(commentId)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	uid, err := GetLoggedUser(r)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError)

		return
	}

	// Insert into Post_Dislike table
	inserterr := database.InsertCommentDislikes(cid, uid)
	if inserterr != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
	}

	database.DeleteCommentLike(cid, uid)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postIDInt), http.StatusSeeOther)
}
