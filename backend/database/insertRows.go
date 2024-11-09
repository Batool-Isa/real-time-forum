package database

import (
	"log"
	"time"
)

func CreateUser(username string, email string, age int, gender string, firstName string, lastName string, pass string) error{
	//preparer statment to create user 
	prepStatment , err := db.Prepare("INSERT INTO User(username, email, age, gender, firstName, lastName, password) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}

	//excute the statment
	_, err = prepStatment.Exec(username, email, age, gender, firstName, lastName, pass)
	if err != nil{
		log.Fatal(err)
	}
	return nil

}

func InsertNewSession(session string, userID int) error{
	//preparer statment to create user 
	prepStatment , err := db.Prepare("INSERT INTO Session(session, timestamp, user_Id) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}

	//excute the statment
	_, err = prepStatment.Exec(session, time.Now().Add(12*time.Hour), userID)
	if err != nil{
		log.Fatal(err)
	}
	return nil

}

