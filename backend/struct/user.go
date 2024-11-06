package structs

import "time"

type Session struct {
	SessionID int       
	Session   string
	Timestamp time.Time 
	UserID    int
}
