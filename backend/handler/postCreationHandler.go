package handler

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/utils"
)

// CreateHandler handles creating a post.
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session from context
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		// If session is nil, user is not authenticated
		fmt.Println("DEBUG: Unauthorized access attempt")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Handle GET request
	if r.Method == "GET" {
		fmt.Println("DEBUG: Rendering create_post.html")
		utils.RenderTemplate(w, r, "create_post.html", session)
		return
	}

	// Handle POST request
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			// Debugging for form parsing errors
			fmt.Println("ERROR: Failed to parse form", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Retrieve form data
		content := r.Form.Get("content")
		category := r.Form["category"]

		// Debugging form data
		fmt.Println("DEBUG: Content:", content)
		fmt.Println("DEBUG: Categories:", category)

		// Retrieve the logged-in user
		uid, err := RetrieveLoggedUser(r)
		if err != nil {
			// Debugging user retrieval
			fmt.Println("ERROR: Unable to retrieve logged-in user", err)
			http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)
			return
		}

		// Insert post into the database
		Inserterr := database.InsertPost(uid, content, category)
		if Inserterr != nil {
			// Debugging database insertion errors
			fmt.Println("ERROR: Failed to insert post", Inserterr)
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		// Redirect to the home page after successful post creation
		fmt.Println("DEBUG: Post created successfully, redirecting to home")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
