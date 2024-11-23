package handler

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/utils"
	"net/http"
	"fmt"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized)
		//http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("http.Redi1354151jrect")
		return
	}
	if r.Method == "GET" {
		utils.RenderTemplate(w, r, "create_post.html", session)
		return
	}
	if r.Method == "POST" {
		r.ParseForm()
		content := r.Form.Get("content")
		category := r.Form["category"]

		uid, err := RetrieveLoggedUser(r)
		if err != nil {
			utils.ErrorHandler(w, r, http.StatusInternalServerError)
			//http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)

			return
		}

		Inserterr := database.InsertPost(uid, content, category)
		if Inserterr != nil {
			utils.ErrorHandler(w, r, http.StatusBadRequest)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)

}
