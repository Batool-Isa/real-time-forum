package handler

import (
	"encoding/json"
//	"fmt"
	"net/http"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"

	"golang.org/x/crypto/bcrypt"
)
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Check for existing session
    session := middleware.GetSessionFromContext(r.Context())
    if session != nil {
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
    password := r.FormValue("password")
    firstName := r.FormValue("firstName")
    lastName := r.FormValue("lastName")
    gender := r.FormValue("gender")

    // Initialize an errors map
    errors := make(map[string]string)

    // Validate required fields
    if username == "" {
        errors["username"] = "Username is required."
    }
    if email == "" {
        errors["email"] = "Email is required."
    }
    if firstName == "" {
        errors["firstName"] = "First name is required."
    }
    if lastName == "" {
        errors["lastName"] = "Last name is required."
    }
    if password == "" {
        errors["password"] = "Password is required."
    }
    if ageStr == "" {
        errors["age"] = "Age is required."
    } else {
        // Convert age to integer
        _, err := strconv.Atoi(ageStr)
        if err != nil {
            errors["age"] = "Age must be a valid number."
        }
    }

    // If there are validation errors, send them as JSON response
    if len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "errors": errors,
        })
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

	  // Convert age from string to int
	  age, err := strconv.Atoi(ageStr)
	  if err != nil {
		  http.Error(w, "Error converting age to integer", http.StatusBadRequest)
		  return
	  }
    // Create the user in the database
    err = database.CreateUser(username, email, age, gender, firstName, lastName, string(hashedPassword))
    if err != nil {
        http.Error(w, "Error registering user", http.StatusInternalServerError)
        return
    }

    // Retrieve the newly created user
    user, err := database.RetrieveUser(username)
    if err != nil {
        http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
        return
    }

    // Create a session for the user
    err = CreateUserSession(w, user.UserID)
    if err != nil {
        http.Error(w, "Error creating session", http.StatusUnauthorized)
        return
    }

    // Send success response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Registration successful.",
    })
}
