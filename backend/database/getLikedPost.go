package database

import (
	"real-time-forum/backend/struct"
	"log"
	"strings"
)

func GetLikedPost(userID int) ([]structs.Post, error) {

	query := `SELECT 
        p.post_id,
        p.user_id,
        p.dislike,
        p.like,
        p.post_heading,
        p.post_data,
        u.username,
        COALESCE(GROUP_CONCAT(c.category_name), ', ', '') AS category_name
    FROM
        posts p 
    INNER JOIN
        users u ON p.user_id = u.uid
    LEFT JOIN
        post_categories pc ON p.post_id = pc.post_id
    LEFT JOIN
        categories c ON pc.category_id = c.category_id
    LEFT JOIN
        likes l ON p.post_id = l.post_id
    WHERE
        l.user_id = ?
    GROUP BY
        p.post_id
    ORDER BY
        p.post_id DESC;
    `

	rows, err := db.Query(query, userID, userID)
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
		post.CategoryName = strings.Split(categoryName, ", ")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
        log.Println(err)
		return nil, err
	}

	return posts, nil
}
