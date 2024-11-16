package handler
import (
    "encoding/json"
    "net/http"
    "real-time-forum/backend/database"
)
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    users, err := database.FetchAllUsers()
    if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
