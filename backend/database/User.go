package database

import (
	"log"
)

func CreateUser(username string, email string, age int, gender string, firstName string, lastName string) error{
	//preparer statment to create user 
	prepStatment , err := db.Prepare("INSERT INTO User(username, email, age, gender, firstName, lastName) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}

	//excute the statment
	_, err = prepStatment.Exec(username, email, age, gender, firstName, lastName)
	if err != nil{
		log.Fatal(err)
	}
	return nil

}

func RetriveSession(sessionCookie string){
	
}