package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
)

type DataPassed struct {
	Username string
	Email    string
	Errors   map[string]string
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check for existing session
	session := middleware.GetSessionFromContext(r.Context())
	if session != nil {
		fmt.Println("ERROR: User is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	// Convert age from string to int
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Error converting age to integer", http.StatusBadRequest)
		fmt.Println("Error converting age to int:", err)
		return
	}

	// Encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		fmt.Println("Error encrypting password:", err)
		return
	}

	// Create the user in the database
	err = database.CreateUser(username, email, age, gender, firstName, lastName, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		fmt.Println("Error creating user in database:", err)
		return
	}

	// Retrieve the newly created user to get the user ID
	user, err := database.RetrieveUser(username)
	if err != nil {
		http.Error(w, "Error fetching user data after registration", http.StatusInternalServerError)
		fmt.Println("Error retrieving user:", err)
		return
	}

	// Create a session for the user
	err = CreateUserSession(w, user.UserID)
	if err != nil {
		http.Error(w, "Error creating user session", http.StatusUnauthorized)
		fmt.Println("Error creating session:", err)
		return
	}

	// Redirect to the homepage on successful registration
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

