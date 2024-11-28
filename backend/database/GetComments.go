package database

import (
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
