package database

import (
	"real-time-forum/backend/struct"
	"database/sql"
	"log"
)

func GetPostCreated(userId string) ([]structs.Post, error) {
	query := `
		SELECT 
			p.post_id, 
			p.post_heading, 
			p.post_data, 
			p.user_id, 
			p.dislike, 
			p.like, 
			c.category_name
		FROM posts p
		JOIN post_categories pc ON p.post_id = pc.post_id
		JOIN categories c ON pc.category_id = c.category_id
		WHERE p.user_id = ?
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var userID sql.NullInt64
		var dislike sql.NullInt64
		var like sql.NullInt64
		var categoryName string

		// Scan NULL-able fields
		err := rows.Scan(&post.PostID, &post.PostHeading, &post.Postdescription, &userID, &dislike, &like, &categoryName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Check for NULL values and convert to appropriate types
		if userID.Valid {
			post.UserID = int(userID.Int64)
		}
		if dislike.Valid {
			post.Dislike = int(dislike.Int64)
		}
		if like.Valid {
			post.Like = int(like.Int64)
		}
		post.CategoryName = []string{categoryName}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return posts, nil
}