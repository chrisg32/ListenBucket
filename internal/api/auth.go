package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/chrisg32/ListenBucket/internal/database"
)

const (
	sessionCookieName = "session"
	sessionDuration   = 30 * 24 * time.Hour // 30 days
)

type contextKey string

const userContextKey contextKey = "user"

// SetupRequest represents the initial setup request
type SetupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response for auth endpoints
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// User represents a user in API responses (without sensitive data)
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// SetupStatusResponse represents the setup status
type SetupStatusResponse struct {
	SetupRequired bool `json:"setup_required"`
}

// getSetupStatus returns whether initial setup is required
// GET /api/auth/setup-status
func (s *Server) getSetupStatus(w http.ResponseWriter, r *http.Request) {
	count, err := s.db.GetUserCount()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	json.NewEncoder(w).Encode(SetupStatusResponse{
		SetupRequired: count == 0,
	})
}

// setup creates the initial admin user
// POST /api/auth/setup
func (s *Server) setup(w http.ResponseWriter, r *http.Request) {
	// Check if setup is already complete
	count, err := s.db.GetUserCount()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}
	if count > 0 {
		writeError(w, http.StatusBadRequest, "setup_complete", "setup has already been completed")
		return
	}

	var req SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if err := validateCredentials(req.Username, req.Password); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "hash_error", "failed to hash password")
		return
	}

	// Create user
	user, err := s.db.CreateUser(req.Username, string(hash))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Create session
	session, err := s.db.CreateSession(user.ID, sessionDuration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Set session cookie
	setSessionCookie(w, session.ID, session.ExpiresAt)

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "setup complete",
		User: &User{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// login authenticates a user
// POST /api/auth/login
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	user, err := s.db.GetUserByUsername(req.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}
	if user == nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid username or password")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid username or password")
		return
	}

	// Create session
	session, err := s.db.CreateSession(user.ID, sessionDuration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Set session cookie
	setSessionCookie(w, session.ID, session.ExpiresAt)

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		User: &User{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// logout invalidates the current session
// POST /api/auth/logout
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		s.db.DeleteSession(cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		Message: "logged out",
	})
}

// me returns the current user
// GET /api/auth/me
func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "not authenticated")
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Success: true,
		User: &User{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// authMiddleware checks for a valid session
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if setup is required
		count, _ := s.db.GetUserCount()
		if count == 0 {
			// No users exist, allow through to show setup screen
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "not authenticated")
			return
		}

		session, err := s.db.GetSession(cookie.Value)
		if err != nil || session == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "invalid or expired session")
			return
		}

		user, err := s.db.GetUserByID(session.UserID)
		if err != nil || user == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "user not found")
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// optionalAuthMiddleware adds user to context if authenticated, but doesn't require it
func (s *Server) optionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err == nil {
			session, _ := s.db.GetSession(cookie.Value)
			if session != nil {
				user, _ := s.db.GetUserByID(session.UserID)
				if user != nil {
					ctx := context.WithValue(r.Context(), userContextKey, user)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func getUserFromContext(ctx context.Context) *database.User {
	user, ok := ctx.Value(userContextKey).(*database.User)
	if !ok {
		return nil
	}
	return user
}

func setSessionCookie(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func validateCredentials(username, password string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return &validationError{"username must be at least 3 characters"}
	}
	if len(password) < 6 {
		return &validationError{"password must be at least 6 characters"}
	}
	return nil
}

type validationError struct {
	message string
}

func (e *validationError) Error() string {
	return e.message
}
