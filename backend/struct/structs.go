package structs

import "time"

type Session struct {
	SessionID int
	Session   string
	Timestamp time.Time
	UserID    int
}

type User struct {
	UserID int
	Username  string
	Email     string
	Age       string
	FirstName string
	LastName  string
	Gender    string
	Password  string
}

type Post struct {
	PostID          int        `json:"postId"`           // Unique identifier for the post
	UserID          int        `json:"userId"`          // User who created the post
	Dislike         int        `json:"dislike"`         // Number of dislikes
	Like            int        `json:"like"`            // Number of likes
	Title           string     `json:"title"`           // Title of the post (was unnamed `string`)
	PostDescription string     `json:"postDescription"` // Description/content of the post
	Username        string     `json:"username"`        // Username of the author
	CategoryName    []string   `json:"categoryName"`    // Categories associated with the post
	Comments        []Comment  `json:"comments"`        // List of comments on the post
}


type Comment struct {
	CommentID int
	UserID    int
	PostID    int
	Text      string
	UserName  string
	CommentLike int
	CommentDislike int
}

type Category struct {
	ID       int
	Category string
}