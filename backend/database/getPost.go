package database

import (
	"real-time-forum/backend/struct"
	"database/sql"
	"log"
	"strings"
)

// var db *sql.DB

func GetAllPosts() ([]structs.Post, error) {
	query := `
	SELECT 
    posts.post_id, 
    posts.user_id, 
    posts.dislike, 
    posts.like, 
    posts.post_heading, 
    posts.post_data, 
    users.username, 
    COALESCE(GROUP_CONCAT(categories.category_name), ', ', '') AS category_name	
FROM 
    posts
INNER JOIN 
    users ON posts.user_id = users.uid
LEFT JOIN 
    post_categories ON posts.post_id = post_categories.post_id
LEFT JOIN 
    categories ON post_categories.category_id = categories.category_id
GROUP BY 
    posts.post_id
ORDER BY
	posts.post_id DESC;`

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
    posts.post_id, 
    posts.user_id, 
    posts.dislike, 
    posts.like, 
    posts.post_heading, 
    posts.post_data, 
    users.username, 
    COALESCE(GROUP_CONCAT(categories.category_name), ', ', '') AS category_name	
FROM 
    posts
INNER JOIN 
    users ON posts.user_id = users.uid
LEFT JOIN 
    post_categories ON posts.post_id = post_categories.post_id
LEFT JOIN 
    categories ON post_categories.category_id = categories.category_id
WHERE
    posts.post_id = ?
GROUP BY 
    posts.post_id;`

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



func GetPostByUserID(userID int) ([]structs.Post, error) {
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
	WHERE
		p.user_id = ?
	GROUP BY
		p.post_id
	ORDER BY
		p.post_id DESC;

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


func GetPostsByCategory(categoryID int) ([]structs.Post, error) {
	
	query := `
		SELECT 
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
