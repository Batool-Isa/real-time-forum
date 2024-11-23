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

// func GetUserIDFromCookie(r *http.Request) (int, error) {
// 	cookie, err := r.Cookie("user_id")
// 	if err != nil {
// 		return 0, err
// 	}

// 	userID, err := strconv.Atoi(cookie.Value)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return userID, nil
// }

func LikePost(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.SessionKey).(structs.Session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		w.WriteHeader(http.StatusForbidden)
		return
		// utils.ErrorHandler(w, r, http.StatusForbidden)

		//http.Error(w, "Unable to retrieve session", http.StatusInternalServerError)
	}
	if r.Method != http.MethodPost {
		utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)

		//http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	if postID == "" {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
		//http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(postID)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
		//http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	uid, err := RetrieveLoggedUser(r)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusForbidden)
		//utils.ErrorHandler(w,r,http.StatusInternalServerError)
		//http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)

		return
	}

	// Insert into Post_Like table
	inserterr := database.InsertLikes(pid, uid)
	if inserterr != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
	}
	database.DeleteDislike(pid, uid)
	database.UpdatePost(pid)

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", pid), http.StatusSeeOther)
}
