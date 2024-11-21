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

func DislikePost(w http.ResponseWriter, r *http.Request) {
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

	postID := r.FormValue("post_id")
	if postID == "" {
		//utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("error")

		//http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(postID)
	if err != nil {
		//utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("error")

		//http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	uid, err := RetrieveLoggedUser(r)
	if err != nil {
	//	utils.ErrorHandler(w, r, http.StatusInternalServerError)
	fmt.Println("error")

		//utils.ErrorHandler(w,r,http.StatusInternalServerError)
		//http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)

		return
	}

	// Insert into dislikes table
	inserterr := database.InsertDislikes(pid, uid)
	if inserterr != nil {
	//	utils.ErrorHandler(w, r, http.StatusBadRequest)
	fmt.Println("error")

	}
	database.DeleteLike(pid, uid)
	database.UpdatePost(pid)

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", pid), http.StatusSeeOther)
}
