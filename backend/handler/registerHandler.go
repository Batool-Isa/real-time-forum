package handler

import(
	"net/http"
	"strconv"
	"fmt"
)


type User struct {
	username string
	email string
	age int
	gender string
	firstName string
	lastName string
}
func registerUser(w http.ResponseWriter, r *http.Request) error{
if (r.Method == http.MethodPost){
		 err := r.ParseForm()
	if err!= nil{
		http.Error(w, "Unable to parse registeration form", http.StatusBadRequest)
		return err
	}

	}	
	//convert age from string to integer
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil{
		http.Error(w, "Invalid age format", http.StatusBadRequest)
		return err
	}

	user := User{
	username : r.FormValue("username"),
	email : r.FormValue("email"),
	age : age,	
	gender : r.FormValue("gender"),
	firstName : r.FormValue("firstName"),
	lastName : r.FormValue("lastName"),
	}
	fmt.Fprintf(w, "User Registered: %+v\n", user)
	return nil




}