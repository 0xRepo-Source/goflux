package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware provides authentication middleware for HTTP handlers
type Middleware struct {
	store *TokenStore
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(store *TokenStore) *Middleware {
	return &Middleware{store: store}
}

// RequireAuth wraps a handler to require authentication
func (m *Middleware) RequireAuth(requiredPermission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		user, permissions, err := m.store.Validate(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		// Check permission
		if requiredPermission != "" && !HasPermission(permissions, requiredPermission) {
			http.Error(w, fmt.Sprintf("Permission denied. Required: %s", requiredPermission), http.StatusForbidden)
			return
		}

		// Set user in request context (optional, for logging)
		r.Header.Set("X-Authenticated-User", user)

		// Call the next handler
		next(w, r)
	}
}

// OptionalAuth wraps a handler to optionally accept authentication
func (m *Middleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				user, _, err := m.store.Validate(parts[1])
				if err == nil {
					r.Header.Set("X-Authenticated-User", user)
				}
			}
		}
		next(w, r)
	}
}
