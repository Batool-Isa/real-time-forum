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
