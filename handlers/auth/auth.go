package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"homeland/config"
	"homeland/models"
	"homeland/utils"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func jsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(db *bun.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			jsonResponse(w, http.StatusBadRequest, ErrorResponse{
				Status:  "error",
				Message: "Invalid request payload",
			})
			return
		}

		var staff models.Staff
		err := db.NewSelect().Model(&staff).Where("email = ?", req.Email).Scan(r.Context())
		if err != nil {
			log.Printf("User not found: %v", err)
			jsonResponse(w, http.StatusUnauthorized, ErrorResponse{
				Status:  "error",
				Message: "Invalid email or password",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(req.Password)); err != nil {
			log.Printf("Incorrect password for user: %s", req.Email)
			jsonResponse(w, http.StatusUnauthorized, ErrorResponse{
				Status:  "error",
				Message: "Your Password is incorrect",
			})
			return
		}

		accessToken, err := utils.GenerateToken(staff.ID, staff.Email, string(staff.Role), cfg.JWTSecret)
		if err != nil {
			log.Printf("Failed to generate access token: %v", err)
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Failed to generate access token",
			})
			return
		}

		refreshToken, err := utils.GenerateRefreshToken(staff.ID, cfg.JWTSecret)
		if err != nil {
			log.Printf("Failed to generate refresh token: %v", err)
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Failed to generate refresh token",
			})
			return
		}

		jsonResponse(w, http.StatusOK, LoginResponse{
			Status:       "success",
			Message:      "Login successful",
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Data: map[string]interface{}{
				"id":               staff.ID,
				"first_name":       staff.FirstName,
				"middle_name":      staff.MiddleName,
				"last_name":        staff.LastName,
				"email":            staff.Email,
				"agent_id":         staff.AgentID,
				"profile_photo":    staff.ProfilePhoto,
				"position":         staff.Position,
				"address":          staff.Address,
				"department":       staff.Department,
				"date_of_birth":    staff.DateOfBirth.Format("2006-01-02"),
				"state_of_origin":  staff.StateOfOrigin,
				"role":             staff.Role,
				"password_changed": staff.MustChangePassword,
			},
		})
	}
}
