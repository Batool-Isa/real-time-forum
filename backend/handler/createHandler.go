
package handler

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
)

// CreateHandler handles creating a post in a single-page application (SPA).
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("ERROR: Failed to parse form", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Debugging raw form data
		fmt.Println("DEBUG: Raw Form Data:", r.Form)

		content := r.Form.Get("content")
		categories := r.Form["category"]

		// Debug parsed values
		fmt.Println("DEBUG: Content:", content)
		fmt.Println("DEBUG: Categories:", categories)

		if content == "" || len(categories) == 0 {
			fmt.Println("ERROR: Content or categories missing")
			http.Error(w, "Content or categories cannot be empty", http.StatusBadRequest)
			return
		}

		uid, err := RetrieveLoggedUser(r)
		if err != nil {
			fmt.Println("ERROR: Unable to retrieve logged-in user", err)
			http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
			return
		}

		err = database.InsertPost(uid, content, categories)
		if err != nil {
			fmt.Println("ERROR: Failed to insert post", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// func PostHandler(w http.ResponseWriter, r *http.Request) {
// 	session := middleware.GetSessionFromContext(r.Context())
// 	postId := r.URL.Query().Get("id")
// 	if postId == "" {
// 		//utils.ErrorHandler(w, r, http.StatusBadRequest)
// 		fmt.Println("ERROR ")
// 		return
// 	}
// 	post_id, err := strconv.Atoi(postId)
// 	if err != nil {
// 		fmt.Println("ERROR CONVERTING")
// 		//utils.ErrorHandler(w, r, http.StatusBadRequest)
// 		return
// 	}
// 	post, err := database.GetPostById(post_id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			//utils.ErrorHandler(w, r, http.StatusNotFound)
// 			fmt.Println("ERROR GETTING POST BY ID")
// 			return
// 		}
// 		//utils.ErrorHandler(w, r, http.StatusInternalServerError)
// 		fmt.Println("ERROR ")
// 		return
// 	}
// 	data := struct {
// 		Post    structs.Post
// 		Session *structs.Session
// 	}{
// 		Post:    post,
// 		Session: session,
// 	}
// 	utils.RenderTemplate(w, r, "comment.html", data)
// }
