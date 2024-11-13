package database

import(
	"database/sql"
	"real-time-forum/backend/struct"
	"log"
	"fmt"
)


func FetchSession(sessionID string)(structs.Session, error){

	//sql query
	query := "SELECT session_Id, session, timestamp, user_Id WHERE session_Id = ?"
	row := db.QueryRow(query, sessionID)
	//declare a session variable
	var userSession structs.Session
	err := row.Scan(&userSession.SessionID, &userSession.Session, &userSession.Timestamp, &userSession.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Session not found:", err)
			return structs.Session{}, fmt.Errorf("session not found")
		}
		log.Println("Error scanning session:", err)
		return structs.Session{}, err
	}
	return userSession, nil

}

func GetSessionByUserID(uID int)(structs.Session, error){

	//sql query       
	query := "SELECT session_Id, session, timestamp, user_Id WHERE user_Id = ?"
	row := db.QueryRow(query, uID)
	//declare a session variable
	var userSession structs.Session
	err := row.Scan(&userSession.SessionID, &userSession.Session, &userSession.Timestamp, &userSession.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Session not found:", err)
			return structs.Session{}, fmt.Errorf("session not found")
		}
		log.Println("Error scanning session:", err)
		return structs.Session{}, err
	}
	return userSession, nil

}