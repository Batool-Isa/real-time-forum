package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"
	"real-time-forum/backend/utils"
)
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    session := middleware.GetSessionFromContext(r.Context())
    if session == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Fetch category filter from query parameters
    category := r.URL.Query().Get("filter-category")
    var posts []structs.Post
    var err error

    // Handle API-specific requests with JSON response
    if r.URL.Path == "/api/posts" {
        if category != "" && category != "all" {
            categoryID, ok := categoriesMap1()[category]
            if !ok {
                http.Error(w, "Invalid category selection", http.StatusBadRequest)
                return
            }
            posts, err = database.GetPostsByCategory(categoryID)
        } else {
            posts, err = database.GetAllPosts()
        }

        if err != nil {
            http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
            return
        }

        // Serve JSON response for API
        response := struct {
            Posts []structs.Post `json:"posts"`
        }{
            Posts: posts,
        }
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        }
        return
    }

    // Serve the HTML template for all other requests
    utils.RenderTemplate(w, r, "index.html", nil)
}


// Map of category names to their IDs
func categoriesMap1() map[string]int {
	return map[string]int{
		"Sports":       1,
		"Technology":   2,
		"Education":    3,
		"Health":       4,
		"Entertainment": 5,
		"Travel":       6,
		"Finance":      7,
		"Culture":      8,
	}
}


func categoriesMap() map[string]int {
	categoryMap := make(map[string]int)
	category, err := database.GetCategories()
	if err != nil {
		fmt.Println("Error getting categories:", err)
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
