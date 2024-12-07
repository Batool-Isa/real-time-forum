package middleware

import (
	"real-time-forum/backend/database"
	"real-time-forum/backend/structs"
	"context"
	"fmt"
	"net/http"
	"time"
)

// sessionContextKey is unexported to prevent collisions
type sessionContextKey int

const (
	SessionKey sessionContextKey = iota
)

// SessionValidator is a middleware function that validates session expiry
func SessionValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session cookie
		sessionCookie, err := r.Cookie("session_Id")
		if err != nil {
			// No session cookie, assume session has expired
			next.ServeHTTP(w, r)
			fmt.Println("Session validation error: no active session cookie")
			return
		}

		// Retrieve session details from the database
		session, err := database.FetchSession(sessionCookie.Value)
		if err != nil {
			// Handle errors if session is not found or expired
			next.ServeHTTP(w, r)
			fmt.Println("Session validation error: unable to fetch session")
			return
		}

		// Verify if the session has expired
		if checkSessionExpiry(session.Session) {
			next.ServeHTTP(w, r)
			fmt.Println("Session validation error: session has expired")
			return
		}

		// Embed session in context
		ctx := context.WithValue(r.Context(), SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// checkSessionExpiry verifies if the given session has expired
func checkSessionExpiry(sessionID string) bool {
	// Retrieve the session data from the session store
	session, err := database.FetchSession(sessionID)
	if err != nil {
		return true // No session available
	}

	// Check if the session has expired based on the timestamp
	return time.Now().After(session.Timestamp)
}

// GetSessionFromContext extracts the session from the context
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


// OptionalSessionMiddleware checks for session and adds it to context if available
func OptionalSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_Id")
		if err != nil {
			next.ServeHTTP(w, r) // No session, proceed without it
			return
		}

		session, err := database.FetchSession(sessionCookie.Value)
		if err != nil || checkSessionExpiry(session.Session) {
			next.ServeHTTP(w, r) // Invalid or expired session, proceed without it
			return
		}

		// Store session in context
		ctx := context.WithValue(r.Context(), SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}