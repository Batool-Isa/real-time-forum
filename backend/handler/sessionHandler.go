package handler

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/struct"
	"real-time-forum/backend/utils"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Global session store and synchronization mechanism
var (
	activeSessions = make(map[string]bool) // Map to store active session IDs
	sessionMutex   = &sync.Mutex{}          // Mutex for thread-safe access to session store
)

// CreateUserSession generates a new session for the user
func CreateUserSession(w http.ResponseWriter, userID int) error {
	// Check if a valid session already exists for the user
	existingSession, err := database.FetchSessionByUser(userID)
	if err == nil && existingSession.Timestamp.After(time.Now()) {
		fmt.Println("User already has a valid session")
		return errors.New("user already has an active session")
	}

	// Generate a unique session ID
	sessionID, err := generateSessionID()
	if err != nil {
		fmt.Println("Error generating session ID:", err)
		return err
	}

	// Store the session in the session map
	sessionMutex.Lock()
	activeSessions[sessionID] = true
	for id := range activeSessions {
		fmt.Println("Current active session ID:", id)
	}
	sessionMutex.Unlock()

	// Send the session ID to the client as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
		// MaxAge: 600000, // Uncomment to set expiration to 10 minutes
	})

	// Save the session in the database
	database.SaveSession(sessionID, userID)
	fmt.Println("Session ID created:", sessionID)
	return nil
}

// generateSessionID creates a unique session ID using random bytes
func generateSessionID() (string, error) {
	// Generate 32 random bytes for the session ID
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert bytes to a hexadecimal string representation
	return hex.EncodeToString(b), nil
}

// RetrieveLoggedUser retrieves the user ID of the logged-in user based on the session cookie
func RetrieveLoggedUser(r *http.Request) (int, error) {
	// Retrieve session cookie from the request
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("No active session cookie found")
		return 0, err
	}

	// Fetch session details using the session ID
	session, err := database.FetchSession(sessionCookie.Value)
	if err != nil {
		fmt.Println("Session retrieval error:", err)
		return 0, err
	}

	return session.UserID, nil
}

// SessionHandler is middleware that verifies session validity before allowing access to protected routes
func SessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session cookie
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println("No active session cookie")
			next.ServeHTTP(w, r)
			return
		}

		// Fetch session details from the database
		session, err := database.FetchSession(sessionCookie.Value)
		if err != nil {
			utils.SendError(w, r, http.StatusUnauthorized)
			fmt.Println("Session error:", err)
			next.ServeHTTP(w, r)
			return
		}

		// Verify session expiration
		if checkSessionExpiration(session.Session) {
			fmt.Println("Session has expired")
			next.ServeHTTP(w, r)
			return
		}

		// Embed session in request context
		ctx := context.WithValue(r.Context(), SessionKey, session)
		fmt.Println("Session is valid and active")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// checkSessionExpiration verifies if the session has expired based on its timestamp
func checkSessionExpiration(sessionID string) bool {
	// Retrieve session data
	session, err := database.FetchSession(sessionID)
	if err != nil {
		return true // Session not found
	}

	// Check if current time is after session timestamp, indicating expiration
	return time.Now().After(session.Timestamp)
}

// GetSessionFromContext extracts the session from the request context
func GetSessionFromContext(ctx context.Context) *structs.Session {
	val := ctx.Value(SessionKey)
	if val == nil {
		return nil
	}
	session, ok := val.(structs.Session)
	if !ok {
		return nil
	}
	return &session
}
