package handler

import(
	"net/http"
	"strconv"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	//"errors"
	"real-time-forum/backend/utils"
	"real-time-forum/backend/database"

)


func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse registration form", http.StatusBadRequest)
		return
	}

	// Retrieve form values
	username := r.FormValue("username")
	email := r.FormValue("email")
	ageStr := r.FormValue("age")
	password := r.FormValue("pass")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	gender := r.FormValue("gender")


	//convert age from string to int 
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		fmt.Print("Error converting age to int")
		return
	}

	//encrypt password
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error encryptong password")
		return 
		//utils.ErrorHandler(w, r, http.StatusInternalServerError)
	}

	// Validate other required fields
	// if err := validateUserInput(username, email, firstName, lastName); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Create the user in the database
	err = database.CreateUser(username, email, age, gender, firstName, lastName, string(pass))
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	// Render login page after successful registration
	utils.RenderTemplate(w, r, "login.html", nil)

}

// // validateUserInput performs basic validation on required fields
// func validateUserInput(username, email, firstName, lastName string)  {
// 	if username == "" || email == "" || firstName == "" || lastName == "" {
// 		fmt.Println("Error Empty Field")
// 		return error
// 	}
// 	// Additional validation checks (like email format) can be added here
// 	return nil
// }

