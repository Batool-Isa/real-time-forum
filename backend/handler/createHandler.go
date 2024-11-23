
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
)

// CreateHandler handles creating a post in a single-page application (SPA).
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session from context
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		// User is not authenticated
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Handle POST request for creating a new post
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			// Log and respond with a JSON error
			fmt.Println("ERROR: Failed to parse form", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid form data"})
			return
		}

		// Retrieve form data
		content := r.Form.Get("content")
		categories := r.Form["category"]

		// Debugging form data
		fmt.Println("DEBUG: Content:", content)
		fmt.Println("DEBUG: Categories:", categories)

		// Retrieve the logged-in user ID
		uid, err := RetrieveLoggedUser(r)
		if err != nil {
			fmt.Println("ERROR: Unable to retrieve logged-in user", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user ID"})
			return
		}

		// Insert post into the database
		Inserterr := database.InsertPost(uid, content, categories)
		if Inserterr != nil {
			fmt.Println("ERROR: Failed to insert post", Inserterr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create post"})
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
		return
	}

	// Method not allowed
	fmt.Println("ERROR: Method not allowed")
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
