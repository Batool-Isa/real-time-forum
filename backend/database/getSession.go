package database

import (
	//"database/sql"
	"log"
	"real-time-forum/backend/structs"
	//"fmt"
)
func FetchSession(sessionID string)(structs.Session, error){
	var session structs.Session
	err := db.QueryRow("SELECT session_Id, session, timestamp, user_Id FROM Session WHERE session = ? ", sessionID).Scan(&session.SessionID, &session.Session, &session.Timestamp ,&session.UserID)
	if err != nil {
		log.Println(err)
		return structs.Session{}, err
	}
	return session, nil
}
func GetSessionByUserID(uID int)(structs.Session, error){
	var session structs.Session
	err := db.QueryRow("SELECT session_Id, session, timestamp, user_Id FROM Session WHERE user_Id = ? ", uID).Scan(&session.SessionID, &session.Session, &session.Timestamp,&session.UserID )
	if err != nil {
		log.Println(err)
		return structs.Session{}, err
	}
	return session, nil
}
