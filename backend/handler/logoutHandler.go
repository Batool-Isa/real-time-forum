package handler

import (
	"real-time-forum/backend/database"
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_Id")
	if err != nil {
		log.Println(err)
	}
	id, err := database.FetchSession(sessionCookie.Value)
	if err != nil {
		log.Println(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_Id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	database.DeleteSession(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
