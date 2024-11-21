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
		fmt.Println("Error")
		//utils.ErrorHandler(w, r, http.StatusForbidden)

		//http.Error(w, "Unable to retrieve session", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		fmt.Println("Error")

//		utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)

		//http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	if postID == "" {
		///utils.ErrorHandler(w, r, http.StatusBadRequest)
		fmt.Println("Error")

		//http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(postID)
	if err != nil {
		fmt.Println("Error")

//		utils.ErrorHandler(w, r, http.StatusBadRequest)
		//http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	uid, err := RetrieveLoggedUser(r)
	if err != nil {
		fmt.Println("Error")
		
		//utils.ErrorHandler(w, r, http.StatusForbidden)
		//utils.ErrorHandler(w,r,http.StatusInternalServerError)
		//http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)

		return
	}

	// Insert into likes table
	inserterr := database.InsertLikes(pid, uid)
	if inserterr != nil {
		fmt.Println("Error")
		
		//utils.ErrorHandler(w, r, http.StatusBadRequest)
	}
	database.DeleteDislike(pid, uid)
	database.UpdatePost(pid)

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", pid), http.StatusSeeOther)
}
