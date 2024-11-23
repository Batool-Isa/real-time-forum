package structs

import (
	"time"
    "github.com/gorilla/websocket"
)

type Session struct {
	SessionID int
	Session   string
	Timestamp time.Time
	UserID    int
	UserName  string
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

type Message struct {
	MessageID  int       `json:"messageId"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	SenderID   int       `json:"senderId"`
	ReceiverID int       `json:"receiverId"`
}


type Client struct {
    Conn     *websocket.Conn
    UserID   int
    Username string
}