package handler
import (
    "encoding/json"
    "net/http"
    "log"
    "real-time-forum/backend/database"
    //"real-time-forum/backend/structs"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    // Retrieve the signed-in user's ID (e.g., from session or token)
    userID, err := RetrieveLoggedUser(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Fetch all users from the database
    allUsers, err := database.FetchAllUsers(userID)
    if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    // Filter out the currently signed-in user
    // var users []structs.UserWithStatus    
    // for _, user := range allUsers {
    //     if user.UserID != userID {
    //         users = append(users, user)
    //     }
    // }
    log.Println(allUsers)
    // Respond with the filtered list of users
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(allUsers)
}
