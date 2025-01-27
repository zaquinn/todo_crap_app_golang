package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

type MiddlewareError struct {
	Error string `json:"error"`
}

func AuthorizationMiddleware(allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		w.Header().Add("Content-Type", "application/json")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(MiddlewareError{Error: "Unauthorized"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userRole := extractRoleFromToken(token)

		if !isRoleAllowed(userRole, allowedRoles) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(MiddlewareError{Error: "Forbidden"})
			return
		}

		next(w, r)
	}
}

func isRoleAllowed(role string, allowedRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}

func extractRoleFromToken(token string) string {
	return token
}
