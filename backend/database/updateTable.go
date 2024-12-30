package database

import (
	// "database/sql"

	"log"
	"real-time-forum/backend/structs"
	"time"
)


func UpdatePost(postId int) (structs.Post, error) {
    // Update likes and dislikes as before
    query := `
        UPDATE Post
        SET 
            likesNum = (SELECT COUNT(*) FROM Post_Like WHERE post_Id = ?),
            dislikesNum = (SELECT COUNT(*) FROM Post_Dislike WHERE post_Id = ?)
        WHERE post_Id = ?;
    `
    stmt, err := db.Prepare(query)
    if err != nil {
        log.Printf("Error preparing update statement: %v", err)
        return structs.Post{}, err
    }
    defer stmt.Close()

    _, err = stmt.Exec(postId, postId, postId)
    if err != nil {
        log.Printf("Error executing update statement: %v", err)
        return structs.Post{}, err
    }

    // Fetch the updated post details
    post, err := GetPostById(postId)
    if err != nil {
        log.Printf("Error fetching updated post: %v", err)
        return structs.Post{}, err
    }

    return post, nil
}

func UpdateSession() error {
	stmt, err := db.Prepare("UPDATE Session SET timestamp = ? WHERE timestamp > ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	newTimestamp := time.Date(2024, 1, 1, 5, 28, 22, 0, time.UTC) // Desired timestamp
	_, err = stmt.Exec(newTimestamp, now)
	return err
}

func GetCommentCount(commentID int) (int, int, error) {
    // Query to fetch likes and dislikes for a specific comment
    query := `
        SELECT 
            (SELECT COUNT(*) FROM Comment_Like WHERE comment_Id = ?) AS likes,
            (SELECT COUNT(*) FROM Comment_Dislike WHERE comment_Id = ?) AS dislikes
    `

    // Prepare the query
    stmt, err := db.Prepare(query)
    if err != nil {
        log.Printf("Error preparing query: %v", err)
        return 0, 0, err
    }
    defer stmt.Close()

    // Execute the query
    var likes, dislikes int
    err = stmt.QueryRow(commentID, commentID).Scan(&likes, &dislikes)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        return 0, 0, err
    }

    // Return likes and dislikes
    return likes, dislikes, nil
}
