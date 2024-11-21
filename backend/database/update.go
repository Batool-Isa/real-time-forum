package database

import (
	// "database/sql"

	"log"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

func UpdatePost(postId int) error {
	stmt1, err := db.Prepare("UPDATE posts SET 'like'= (SELECT count(*) FROM likes where post_id=?) where post_id=?;")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(postId, postId)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt2, err := db.Prepare("UPDATE posts SET dislike=(SELECT count(*) FROM dislikes where post_id=?) where post_id=?;")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt2.Exec(postId, postId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UpdateSession() error {
    stmt, err := db.Prepare("UPDATE sessions SET timestamp = ? WHERE timestamp > ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    now := time.Now()
    newTimestamp := time.Date(2024, 1, 1, 5, 28, 22, 0, time.UTC) // Desired timestamp
    _, err = stmt.Exec(newTimestamp, now)
    return err
}