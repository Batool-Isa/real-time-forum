package database

import (
	"real-time-forum/backend/struct"
	"log"
)

func GetAllComments(postID int) ([]structs.Comment, error) {
	query := `
    SELECT
        comments.comment_id,
        comments.user_id,
        comments.post_id,
        comments.comment,
        users.username,
        (SELECT COUNT(*) FROM likeComment WHERE likeComment.comment_id = comments.comment_id) AS like_count,
        (SELECT COUNT(*) FROM dislikeComment WHERE dislikeComment.comment_id = comments.comment_id) AS dislike_count
    FROM
        comments
    INNER JOIN
        users ON comments.user_id = users.uid
    WHERE
        comments.post_id = ?`

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
		var likeComment int
		err := rows.Scan(&comm.CommentID, &comm.UserID, &comm.PostID, &comm.Text, &comm.UserName, &likeComment, &dislikeComment)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		comm.CommentLike = likeComment
		comm.CommentDislike = dislikeComment
		comms = append(comms, comm)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return comms, nil
}
