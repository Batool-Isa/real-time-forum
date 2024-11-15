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
	PostID          int
	UserID          int
	Dislike         int
	Like            int
	PostHeading     string
	Postdescription string
	Username        string
	CategoryName    []string
	Comments        []Comment
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