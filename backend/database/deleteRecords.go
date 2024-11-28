package database

import (
	// "database/sql"
	"log"
	"real-time-forum/backend/structs"
)

func DeleteLike(postId int, userId int) error {
	stmt1, err := db.Prepare("DELETE FROM Post_Like WHERE post_Id = ? AND user_Id= ?")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(postId, userId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteDislike(postId int, userId int) error {
	stmt1, err := db.Prepare("DELETE FROM Post_Dislike WHERE post_Id = ? AND user_Id= ?")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(postId, userId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteCommentLike(CommentID int, userId int) error {
	stmt1, err := db.Prepare("DELETE FROM Comment_Like WHERE comment_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(CommentID, userId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteCommentDislike(CommentID int, userId int) error {
	stmt1, err := db.Prepare("DELETE FROM Comment_Dislike WHERE comment_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(CommentID, userId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteCategory(catId int) error {
	stmt1, err := db.Prepare("DELETE FROM Category WHERE category_id = ? ")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(catId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func DeleteSession(session structs.Session) error {
	stmt1, err := db.Prepare("DELETE FROM Session WHERE session_id = ? ")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(session.SessionID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func CleanUpPosts() error {
	stmt, err := db.Prepare(`
        DELETE FROM Post
        WHERE post_id NOT IN (
            SELECT post_id
            FROM Post_Category)`,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
