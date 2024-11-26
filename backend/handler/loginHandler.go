package handler

import (
	"encoding/json"
	"fmt"

	//	"fmt"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/structs"
	"strings"

	"golang.org/x/crypto/bcrypt"
)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Check if the user is already logged in
    session := middleware.GetSessionFromContext(r.Context())
    if session != nil {
        fmt.Println("User is already logged in")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Parse form data
    r.ParseForm()
    usernameOrEmail := strings.TrimSpace(r.Form.Get("usernameOrEmail"))
    password := strings.TrimSpace(r.Form.Get("password"))

    fmt.Println("Raw form data:", r.Form)                  // Debugging output
    fmt.Println("Parsed usernameOrEmail:", usernameOrEmail) // Debugging output
    fmt.Println("Parsed password:", password)              // Debugging output

    // Initialize errors map
    errors := make(map[string]string)

    if usernameOrEmail == "" || password == "" {
        errors["user"] = "Username/Email and Password are required"
    }

    // Retrieve user from the database
    var user structs.User
        var err error
        user, err = database.RetrieveUser(usernameOrEmail)
        if err != nil {
            errors["user"] = "Invalid username/email or password"
            fmt.Println("Error retrieving user:", err) // Debugging output
        } else {
            fmt.Printf("User retrieved: %+v\n", user) // Debugging output
        }
    

    // Compare passwords if user is found
        fmt.Println("Entered password:", password)
        fmt.Println("Stored hashed password:", user.Password)
        if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
            errors["user"] = "Invalid username/email or password"
            fmt.Println("Password comparison failed")
        } else {
            fmt.Println("Password comparison successful")
        }
    

    // If there are errors, return them as a JSON response
    if len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "errors": errors,
        })
        return
    }
    if len(errors) > 0 {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "errors": errors,
    })
    return
}


    // Create user session on successful login
    err = CreateUserSession(w, user.UserID)
    if err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        fmt.Println("Error creating session:", err)
        return
    }

    // w.Header().Set("Content-Type", "application/json")
    // json.NewEncoder(w).Encode(map[string]interface{}{
    //     "message": "Login successful",
    // })
    http.Redirect(w, r, "/api/posts", http.StatusSeeOther)

}


// func LoginHandler(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodPost {
//         http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
//         return
//     }
//     // Check if the user is already logged in
//     session := middleware.GetSessionFromContext(r.Context())
//     if session != nil {
//         http.Redirect(w, r, "/", http.StatusSeeOther)
//         return
//     }
//     // Parse form data
//     r.ParseForm()
//     usernameOrEmail := strings.TrimSpace(r.Form.Get("usernameOrEmail"))
//     password := strings.TrimSpace(r.Form.Get("password"))
//     // Initialize errors map
//     errors := make(map[string]string)
//     // Retrieve user from the database
//     user, err := database.RetrieveUser(usernameOrEmail)
//     if err != nil {
//         errors["user"] = "Invalid username/email or password"
//         fmt.Println("Error retrieving user:", err)
//     } else {
//         // Debugging: Print out entered password and stored hash
//         fmt.Println("Entered password:", password)
//         fmt.Println("Stored password hash:", user.Password)
//         // Compare the entered password with the hashed password
//         err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//         if err != nil {
//             errors["user"] = "Invalid username/email or password"
//             fmt.Println("Password comparison failed:", err)
//         }
//     }
//     // If there are errors, render the login page again with errors
//     if len(errors) > 0 {
//         fmt.Println("Error during login")
//         return
//     }
//     // Create user session on successful login
//     err1 := CreateUserSession(w, user.UserID)
//     if err1 != nil {
//         fmt.Println("Error creating session:", err1)
//         return
//     }
//     // Redirect to the homepage on successful login
//     http.Redirect(w, r, "/", http.StatusSeeOther)
// }
