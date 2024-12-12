package database

import (
	"database/sql"
	"log"
	"real-time-forum/backend/structs"
)

func GetAllComments(postID int) ([]structs.Comment, error) {
	query := `
    SELECT
        Comment.comment_Id,
        Comment.user_Id,
        Comment.post_Id,
        Comment.comment_content,
        User.username,
        (SELECT COUNT(*) FROM Comment_Like WHERE Comment_Like.comment_Id = Comment.comment_Id) AS like_count,
        (SELECT COUNT(*) FROM Comment_Dislike WHERE Comment_Dislike.comment_Id = Comment.comment_Id) AS dislike_count
    FROM
		Comment
    INNER JOIN
        User ON Comment.user_Id = User.user_Id
    WHERE
         Comment.post_Id = ?`

	rows, err := db.Query(query, postID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var comms []structs.Comment
	for rows.Next() {
		var comm structs.Comment
		var dislikeComment int
		var Comment_Like int
		err := rows.Scan(&comm.CommentID, &comm.UserID, &comm.PostID, &comm.Text, &comm.UserName, &Comment_Like, &dislikeComment)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		comm.CommentLike = Comment_Like
		comm.CommentDislike = dislikeComment
		comms = append(comms, comm)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return comms, nil
}



// InsertLikeForComment inserts a like for a specific comment by a user.
func InsertLikeForComment(commentID int, userID int) error {
	query := `
	INSERT OR IGNORE INTO Comment_Like (comment_Id, user_Id)
	VALUES (?, ?);`

	_, err := db.Exec(query, commentID, userID)
	if err != nil {
		log.Printf("Error inserting like for comment %d by user %d: %v", commentID, userID, err)
		return err
	}
	return nil
}

// InsertDislikeForComment inserts a dislike for a specific comment by a user.
func InsertDislikeForComment(commentID int, userID int) error {
	query := `
	INSERT OR IGNORE INTO Comment_Dislike (comment_Id, user_Id)
	VALUES (?, ?);`

	_, err := db.Exec(query, commentID, userID)
	if err != nil {
		log.Printf("Error inserting dislike for comment %d by user %d: %v", commentID, userID, err)
		return err
	}
	return nil
}

// GetCommentByID retrieves a comment by its ID.
func GetCommentByID(commentID int) (map[string]interface{}, error) {
	query := `
	SELECT c.comment_Id, c.comment_content, c.user_Id, c.post_Id,
	       (SELECT COUNT(*) FROM Comment_Like WHERE comment_Id = c.comment_Id) AS like_count,
	       (SELECT COUNT(*) FROM Comment_Dislike WHERE comment_Id = c.comment_Id) AS dislike_count
	FROM Comment c
	WHERE c.comment_Id = ?;`

	row := db.QueryRow(query, commentID)

	var commentIDResult int
	var commentContent string
	var userID int
	var postID int
	var likeCount int
	var dislikeCount int

	err := row.Scan(&commentIDResult, &commentContent, &userID, &postID, &likeCount, &dislikeCount)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No comment found with ID %d", commentID)
			return nil, nil
		}
		log.Printf("Error retrieving comment with ID %d: %v", commentID, err)
		return nil, err
	}

	comment := map[string]interface{}{
		"commentID":   commentIDResult,
		"content":     commentContent,
		"userID":      userID,
		"postID":      postID,
		"likeCount":   likeCount,
		"dislikeCount": dislikeCount,
	}

	return comment, nil
}
