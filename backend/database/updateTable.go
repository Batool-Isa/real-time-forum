package database

import (
	// "database/sql"

	"log"
	"time"

)

func UpdatePost(postId int) error {
	stmt1, err := db.Prepare("UPDATE Post SET 'like'= (SELECT count(*) FROM Post_Like where post_id=?) where post_id=?;")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt1.Exec(postId, postId)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt2, err := db.Prepare("UPDATE Post SET dislike=(SELECT count(*) FROM Post_Dislike where post_id=?) where post_id=?;")
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
