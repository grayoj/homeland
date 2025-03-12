package auth

import (
	"net/http"
	"strings"

	"homeland/config"
	"homeland/models"
	"homeland/utils"

	"github.com/uptrace/bun"
)

type AuthResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func AuthCheckHandler(db *bun.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonResponse(w, http.StatusUnauthorized, AuthResponse{
				Status:  "error",
				Message: "Authorization token is required",
			})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			jsonResponse(w, http.StatusUnauthorized, AuthResponse{
				Status:  "error",
				Message: "Invalid authorization format",
			})
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], cfg.JWTSecret)
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, AuthResponse{
				Status:  "error",
				Message: "Invalid or expired token",
			})
			return
		}

		userID := claims.UserID

		var staff models.Staff
		err = db.NewSelect().Model(&staff).Where("id = ?", userID).Scan(r.Context())
		if err != nil {
			jsonResponse(w, http.StatusNotFound, AuthResponse{
				Status:  "error",
				Message: "User not found",
			})
			return
		}

		jsonResponse(w, http.StatusOK, AuthResponse{
			Status:  "success",
			Message: "User authenticated",
			Data: map[string]interface{}{
				"id":               staff.ID,
				"email":            staff.Email,
				"agent_id":         staff.AgentID,
				"role":             staff.Role,
				"password_changed": staff.MustChangePassword,
			},
		})
	}
}
