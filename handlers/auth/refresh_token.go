package auth

import (
	"encoding/json"
	"net/http"
	"strconv"

	"homeland/config"
	"homeland/utils"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

func RefreshTokenHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"status":"error","message":"Invalid payload"}`, http.StatusBadRequest)
			return
		}

		claims, err := utils.ValidateRefreshToken(req.RefreshToken, cfg.JWTSecret)
		if err != nil {
			http.Error(w, `{"status":"error","message":"Invalid or expired refresh token"}`, http.StatusUnauthorized)
			return
		}

		userID, _ := strconv.ParseInt(claims.Subject, 10, 64)
		newAccessToken, err := utils.GenerateToken(userID, "", "", cfg.JWTSecret)
		if err != nil {
			http.Error(w, `{"status":"error","message":"Could not generate new access token"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RefreshResponse{
			Status:      "success",
			Message:     "New access token generated",
			AccessToken: newAccessToken,
		})
	}
}
