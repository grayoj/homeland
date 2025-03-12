package middleware

import (
	"context"
	"net/http"
	"strings"

	"homeland/utils"
)

type contextKey string

const (
	ContextKeyClaims = contextKey("claims")
)

func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			claims, err := utils.ValidateToken(parts[1], jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleAuthorization(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ContextKeyClaims).(*utils.Claims)
			if !ok {
				http.Error(w, "No claims found", http.StatusForbidden)
				return
			}
			if _, allowed := roleSet[claims.Role]; !allowed {
				http.Error(w, "Insufficient privileges", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func FromContext(r *http.Request) *utils.Claims {
	claims, _ := r.Context().Value(ContextKeyClaims).(*utils.Claims)
	return claims
}
