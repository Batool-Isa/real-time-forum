package database

import (
	"database/sql"
	"fmt"
	"log"
	"real-time-forum/backend/structs"
)

// function to get user info from database using his email or username
func RetrieveUser(username string)(structs.User, error){
    //sql query
	query := "SELECT user_Id, username, email, age, gender, firstName, lastName, password FROM User WHERE username = ? OR email = ?"
    row := db.QueryRow(query, username, username)
    // Declare a user variable
    var user structs.User
    err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("User not found:", err)
            return structs.User{}, fmt.Errorf("user not found")
        }
        log.Println("Error scanning user:", err)
        return structs.User{}, err
    }
    return user, nil

}
func GetUsername(uid int) (string, error) {
    // SQL query to get the username based on the user ID
    query := "SELECT username FROM User WHERE user_Id = ?"
    row := db.QueryRow(query, uid)

    // Declare a variable to store the username
    var username string
    err := row.Scan(&username)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("No user found with the given ID:", uid)
            return "", err
        }
        log.Println("Error scanning username:", err)
        return "", err
    }
    return username, nil
}

func FetchAllUsers() ([]structs.User, error) {
    query := "SELECT user_Id, username, email, age, gender, firstName, lastName FROM User"
    rows, err := db.Query(query)
    if err != nil {
        log.Println("Error querying users:", err)
        return nil, err
    }
    defer rows.Close()

    var users []structs.User
    for rows.Next() {
        var user structs.User
        err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Age, &user.Gender, &user.FirstName, &user.LastName)
        if err != nil {
            log.Println("Error scanning user:", err)
            continue
        }
        users = append(users, user)
    }

    return users, nil
}
