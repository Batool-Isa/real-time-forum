package database

import (
	// "database/sql"
	"database/sql"
	"log"
	"real-time-forum/backend/structs"
	"strings"
)

// var db *sql.DB

func GetAllPosts() ([]structs.Post, error) {
	query := `
	SELECT 
    Post.post_id, 
    Post.user_id, 
    Post.dislike, 
    Post.likes, 
    Post.post_heading, 
    Post.post_content, 
    User.username, 
    COALESCE(GROUP_CONCAT(Category.category_name), ', ', '') AS category_name	
FROM 
    Post
INNER JOIN 
    User ON Post.user_id = User.user_id
LEFT JOIN 
    post_Category ON Post.post_id = Post_Category.post_id
LEFT JOIN 
    Category ON post_Category.category_id = Category.category_id
GROUP BY 
    Post.post_id
ORDER BY
	Post.post_id DESC;`

	rows, err := db.Query(query)
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

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return posts, nil
}

func GetPostById(id int) (structs.Post, error) {
	query := `
	SELECT 
    Post.post_id, 
    Post.user_id, 
    Post.dislike, 
    Post.likes, 
    Post.post_heading, 
    Post.post_content, 
    User.username, 
    COALESCE(GROUP_CONCAT(Category.category_name), ', ', '') AS category_name	
FROM 
    Post
INNER JOIN 
    User ON Post.user_id = User.user_id
LEFT JOIN 
    Post_Category ON Post.post_id = Post_Category.post_id
LEFT JOIN 
    categories ON post_Category.category_id = Category.category_id
WHERE
    Post.post_id = ?
GROUP BY 
    Post.post_id;`

	row := db.QueryRow(query, id)

	var post structs.Post
	var categoryName string
	err := row.Scan(&post.PostID, &post.UserID, &post.Dislike, &post.Like, &post.PostHeading, &post.Postdescription, &post.Username, &categoryName)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return structs.Post{}, err // No post found with the given ID
		}
		return structs.Post{}, err
	}
	post.CategoryName = strings.Split(categoryName, ", ")
	post.Comments, err = GetAllComments(id)
	if err != nil {
		log.Println(err)
		return structs.Post{}, err
	}
	return post, nil
}
