package handler

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/utils"
	"fmt"
	"net/http"
	
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		//utils.ErrorHandler(w, r, http.StatusUnauthorized)
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
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		category := r.Form["category"]

		uid, err := RetrieveLoggedUser(r)
		if err != nil {
			//utils.ErrorHandler(w, r, http.StatusInternalServerError)
			//http.Error(w, "Unable to retrieve user ID", http.StatusInternalServerError)

			return
		}

		Inserterr := database.InsertPost(uid, title, content, category)
		if Inserterr != nil {
			//utils.ErrorHandler(w, r, http.StatusBadRequest)
			fmt.Println("ERROR INSERTING")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	fmt.Println("ERROR ")
	//utils.ErrorHandler(w, r, http.StatusMethodNotAllowed)

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
