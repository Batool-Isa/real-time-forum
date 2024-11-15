package database

import (
	"log"
	"real-time-forum/backend/structs"
	"strings"
)

func GetPostsByCategory(categoryID int) ([]structs.Post, error) {

	query := `
		SELECT 
			p.post_id, 
			p.user_id, 
			p.dislike, 
			p.likes, 
			p.post_heading, 
			p.post_content, 
			u.username, 
			COALESCE(GROUP_CONCAT(c.category_name), ', ', '') AS category_name
		FROM 
			Post p
		INNER JOIN 
			User u ON p.user_id = u.user_id
		LEFT JOIN 
			Post_Category pc ON p.post_id = pc.post_id
		LEFT JOIN 
			Category c ON pc.category_id = c.category_id
		WHERE 
			pc.category_id = ?
		GROUP BY 
			p.post_id
		ORDER BY
			p.post_id DESC;
	`

	rows, err := db.Query(query, categoryID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var categoryName string
		err := rows.Scan(&post.PostID, &post.UserID, &post.Dislike, &post.Like, &post.PostHeading, &post.Postdescription, &post.Username, &categoryName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		post.CategoryName = strings.Split(categoryName, ",")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return posts, nil
}
