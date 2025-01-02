package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"regexp"
	"strconv"
	"unicode"

	"golang.org/x/crypto/bcrypt"
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
    log.Println("Username:", username)
    email := r.FormValue("email")
    log.Println("Email:", email)
    ageStr := r.FormValue("age")
        log.Println("Age:", ageStr)
    password := r.FormValue("password")

    confirmPassword := r.FormValue("confirmPassword")
    firstName := r.FormValue("firstName")
        log.Println("First Name:", firstName)
    lastName := r.FormValue("lastName")
    log.Println("Last Name:", lastName)
    // Initialize an errors map
    errors := make(map[string]string)

    // Validate Username
    if len(username) < 3 {
        errors["username"] = "Username must be at least 3 characters long."
    }

    // Validate Email
    emailRegex := `^[^\s@]+@[^\s@]+\.[^\s@]+`
    matched, _ := regexp.MatchString(emailRegex, email)
    if !matched {
        errors["email"] = "Invalid email format."
    }

    // Validate First and Last Name
    if len(firstName) < 3 {
        errors["firstName"] = "First name must be at least 3 characters long."
    }
    if len(lastName) < 3 {
        errors["lastName"] = "Last name must be at least 3 characters long."
    }

    // Validate Age
    age, err := strconv.Atoi(ageStr)
    if err != nil || age <= 0 {
        errors["age"] = "Age must be a valid positive number."
    }

    // Validate Password
    if len(password) < 8 || !containsLetterAndNumber(password) {
        errors["password-p"] = "Password must be at least 8 characters long and include letters and numbers."
    }

    // Validate Confirm Password
    if password != confirmPassword {
        errors["confirmPassword"] = "Passwords do not match."
    }

    // Debug log all errors
    log.Printf("Validation errors: %+v\n", errors)

    // Return errors if validation fails
    if len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
        return
    }
    
    // Proceed with user registration logic (e.g., hashing password, saving to DB)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    err = database.CreateUser(username, email, age, "", firstName, lastName, string(hashedPassword))
    if err != nil {
        http.Error(w, "Error registering user", http.StatusInternalServerError)
        return
    }


    // Redirect on successful registration
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func containsLetterAndNumber(s string) bool {
    hasLetter := false
    hasNumber := false
    for _, char := range s {
        if unicode.IsLetter(char) {
            hasLetter = true
        } else if unicode.IsNumber(char) {
            hasNumber = true
        }
        if hasLetter && hasNumber {
            return true
        }
    }
    return false
}
