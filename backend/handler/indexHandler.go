package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/structs"
	"real-time-forum/backend/utils"
	"real-time-forum/backend/middleware"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "1212121212",
		Path:   "/",
		MaxAge: 600000, // Set the expiration time to 600 seconds (10 minute)
	})
	var categoryMap = categoriesMap()

	if r.URL.Path != "/" {
		utils.ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	posts, err := database.GetAllPosts()
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	categorySelection := r.URL.Query().Get("filter-category")
	if categorySelection != "" && categorySelection != "all" {
		categoryID, ok := categoryMap[categorySelection]

		if !ok {
			// Handle invalid category selection
			utils.ErrorHandler(w, r, http.StatusBadRequest)
			return
		}
		posts, err = database.GetPostsByCategory(categoryID)
		if err != nil {
			utils.ErrorHandler(w, r, http.StatusInternalServerError)

			return
		}
	}

	session := middleware.GetSessionFromContext(r.Context())

	myPosts := r.URL.Query().Get("filter-mypost")
	likedepost := r.URL.Query().Get("filter-likedpost")

	if session != nil {
		userId := session.UserID
		if myPosts == "true" {
			posts, err = database.GetPostByUserID(userId)
			if err != nil {
				utils.ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
		} else if likedepost == "true" {
			posts, err = database.GetLikedPost(userId)
			if err != nil {
				utils.ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}

		}
	}

	//  posts = database.GetPostByUserID(userId)
	category, err := database.GetCategories()
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest)
	}

	data := struct {
		Posts    []structs.Post
		Session  *structs.Session
		Category []structs.Category
	}{
		Posts:    posts,
		Session:  session,
		Category: category,
	}

	utils.RenderTemplate(w, r, "index.html", data)

}

func categoriesMap() map[string]int {
	categoryMap := make(map[string]int)
	category, err := database.GetCategories()
	if err != nil {
		fmt.Println("Error getting categories")
		utils.ErrorHandler(nil, nil, http.StatusBadRequest)
		return nil
	}
	for i, cat := range category {
		cat.ID = i + 1
		
		categoryMap[cat.Category] = cat.ID
	}
	return categoryMap
}

// New SPAHandler
func SPAHandler(w http.ResponseWriter, r *http.Request) {
    // Serve index.html for all routes
    http.ServeFile(w, r, "template/index.html")
}

// New GetPostsHandler
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
    posts, err := database.GetAllPosts()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

// New GetCategoriesHandler
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
    categories, err := database.GetCategories()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}
