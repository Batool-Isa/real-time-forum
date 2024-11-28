package database

import (
	"database/sql"
	"log"
	"real-time-forum/backend/structs"
	"strings"
)

// var db *sql.DB
func GetAllPosts() ([]structs.Post, error) {
	query := `
	SELECT 
    Post.post_Id, 
	Post.post_content, 
	Post.user_Id,
     Post.likesNum, 
    Post.dislikesNum,  
    User.username, 
    COALESCE(GROUP_CONCAT(Category.category_name), ', ', '') AS category_name	
FROM 
    Post
INNER JOIN 
    User ON Post.user_Id = User.user_Id
LEFT JOIN 
    Post_Category ON Post.post_Id = Post_Category.post_Id
LEFT JOIN 
    Category ON Post_Category.category_Id = Category.category_Id
GROUP BY 
    Post.post_Id
ORDER BY
	Post.post_Id DESC;`
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
		err := rows.Scan(&post.PostID, &post.PostDescription, &post.UserID, &post.Like, &post.Dislike, &post.Username, &categoryName)
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
  Post.post_Id, 
	 Post.post_content, 
    Post.user_Id,
	  Post.likesNum, 
    Post.dislikesNum,  
    User.username, 
    COALESCE(GROUP_CONCAT(Category.category_name), ', ', '') AS category_name	
FROM 
    Post
INNER JOIN 
    User ON Post.user_Id = User.user_Id
LEFT JOIN 
    Post_Category ON Post.post_Id = Post_Category.post_Id
LEFT JOIN 
    Category ON Post_Category.category_Id = Category.category_Id
WHERE
    Post.post_Id = ?
GROUP BY 
    Post.post_Id;`
	row := db.QueryRow(query, id)
	var post structs.Post
	var categoryName string
	err := row.Scan(&post.PostID, &post.PostDescription, &post.UserID,&post.Like, &post.Dislike, &post.Username, &categoryName)
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
func GetPostByUserID(userID int) ([]structs.Post, error) {
	query := `SELECT 
	Post.post_Id, 
	Post.post_content, 
    Post.user_Id,
	Post.likesNum, 
	Post.dislikesNum,  
	User.username, 
	  COALESCE(GROUP_CONCAT(Category.category_name), ', ', '') AS category_name	
	FROM
		Post p
	INNER JOIN
		User u ON p.user_Id = u.User_Id
	LEFT JOIN
		Post_Category pc ON p.post_Id = pc.post_Id
	LEFT JOIN
		Category c ON pc.category_Id = c.category_Id
	WHERE
		p.user_Id = ?
	GROUP BY
		p.post_Id
	ORDER BY
		p.post_Id DESC;
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var categoryName string
		err := rows.Scan(&post.PostID, &post.PostDescription, &post.UserID, &post.Like, &post.Dislike, &post.Username, &categoryName)
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
func GetPostsByCategory(categoryID int) ([]structs.Post, error) {

	query := ` 
	SELECT 
	Post.post_Id, 
	Post.post_content, 
	Post.user_Id,
	Post.likesNum, 
	Post.dislikesNum,  
	User.username, 
			COALESCE(GROUP_CONCAT(c.category_name), ', ', '') AS category_name
		FROM 
			Post p
		INNER JOIN 
			User u ON p.user_Id = u.user_Id
		LEFT JOIN 
			Post_Category pc ON p.post_Id = pc.post_Id
		LEFT JOIN 
			Category c ON pc.category_Id = c.category_Id
		WHERE 
			pc.category_Id = ?
		GROUP BY 
			p.post_Id
		ORDER BY
			p.post_Id DESC;
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
		err := rows.Scan(&post.PostID, &post.PostDescription,&post.UserID,&post.Like, &post.Dislike,&post.Username, &categoryName)
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

// GetAllPostsWithPagination fetches all posts with limit and offset
func GetAllPostsWithPagination(limit, offset int) ([]structs.Post, error) {
	rows, err := db.Query("SELECT * FROM Post ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.PostID, &post.PostDescription, &post.Username); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetTotalPostCount fetches the total number of posts
func GetTotalPostCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Post").Scan(&count)
	return count, err
}

func GetPostsByCategoryWithPagination(categoryID, limit, offset int) ([]structs.Post, error) {
	// Prepare the query to fetch posts by category with pagination
	query := `
		SELECT p.id, p.post_content, p.created_at, p.user_id, u.username
		FROM Post p
		INNER JOIN Post_Category pc ON p.id = pc.post_id
		INNER JOIN User u ON p.user_id = u.id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.PostID, &post.PostDescription, &post.UserID, &post.Username)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
