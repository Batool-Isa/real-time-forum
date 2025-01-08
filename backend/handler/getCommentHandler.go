package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"strconv"
	"strings"
)


func CommentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
        return
    }

    // Retrieve the logged-in user's ID
    uid, err := RetrieveLoggedUser(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    // Parse the JSON body
    var reqData struct {
        PostID      int    `json:"post_id"`
        CommentText string `json:"comment"`
    }

    err = json.NewDecoder(r.Body).Decode(&reqData)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    // Validate input
    if reqData.PostID == 0 || reqData.CommentText == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Missing post ID or comment text"})
        return
    }

    // Insert the comment into the database and get the comment ID
    commentID, err := database.InsertComment(reqData.CommentText, uid, reqData.PostID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add comment"})
        return
    }

    // Return success response with the comment ID
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":    "Comment added successfully",
        "comment_id": commentID,
    })
}




// Like a comment
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body
	var requestBody struct {
		CommentID int `json:"comment_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.CommentID <= 0 {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Retrieve the logged-in user's ID
	userID, err := RetrieveLoggedUser(r)
	if userID == 0 || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Process like logic
	err = database.InsertLikeForComment(requestBody.CommentID, userID)
	if err != nil {
		http.Error(w, "Failed to like comment", http.StatusInternalServerError)
		return
	}
database.DeleteCommentDislike(requestBody.CommentID, userID)
	
    // Get updated like and dislike counts
    likeCount, dislikeCount, err := database.GetCommentCount(requestBody.CommentID)
    if err != nil {
        http.Error(w, "Failed to retrieve comment count", http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "commentId": requestBody.CommentID,
        "likes":     likeCount,
        "dislikes":  dislikeCount,
    }
  
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}


func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var requestBody struct {
        CommentID int `json:"comment_id"`
    }

    // Decode the JSON body
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil || requestBody.CommentID <= 0 {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    // Log received comment ID
    log.Printf("Received comment ID: %d", requestBody.CommentID)

    // Retrieve logged-in user
    userID, err := RetrieveLoggedUser(r)
    if userID == 0 || err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Process dislike logic
    err = database.InsertDislikeForComment(requestBody.CommentID, userID)
    if err != nil {
        http.Error(w, "Failed to dislike comment", http.StatusInternalServerError)
        return
    }

    // Optionally remove the like if it exists
    err = database.DeleteCommentLike(requestBody.CommentID, userID)
    if err != nil {
        log.Printf("Failed to remove like for comment ID %d: %v", requestBody.CommentID, err)
    }

    // Get updated like and dislike counts
    likeCount, dislikeCount, err := database.GetCommentCount(requestBody.CommentID)
    if err != nil {
        http.Error(w, "Failed to retrieve comment count", http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "commentId": requestBody.CommentID,
        "likes":     likeCount,
        "dislikes":  dislikeCount,
    }
  
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}



// Get comment details
func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the comment ID from the URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(segments[3])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}
