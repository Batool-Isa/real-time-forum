package main

import (
	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	_ "github.com/mattn/go-sqlite3"
	"log"
	
)

func main(){

	err:= database.CreateDB("forum.db")
	if err != nil {
		log.Fatalf("Error creating or connecting to database: %v", err)
	}
	//create tables
	database.CreateTables()
	
	fmt.Println("Databse Setup complete")	
	fmt.Println("Server Started at http://localhost:8888/")

	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("Error starting server at 8888",err)
	}
}