package database

import (
	// "database/sql"
	"real-time-forum/backend/struct"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DeleteLike(postId int, userId int) error {
	stmt1, err := db.Prepare("DELETE FROM likes WHERE post_id = ? AND user_id= ?")
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
	stmt1, err := db.Prepare("DELETE FROM dislikes WHERE post_id = ? AND user_id= ?")
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
	stmt1, err := db.Prepare("DELETE FROM likeComment WHERE comment_id = ? AND user_id = ?")
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
	stmt1, err := db.Prepare("DELETE FROM dislikeComment WHERE comment_id = ? AND user_id = ?")
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
	stmt1, err := db.Prepare("DELETE FROM categories WHERE category_id = ? ")
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
	stmt1, err := db.Prepare("DELETE FROM sessions WHERE session_id = ? ")
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
        DELETE FROM posts
        WHERE post_id NOT IN (
            SELECT post_id
            FROM post_categories)`,
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
