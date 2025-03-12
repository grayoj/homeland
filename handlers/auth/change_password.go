package auth

import (
	"encoding/json"
	"net/http"

	"homeland/config"
	"homeland/models"
	"homeland/utils"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePasswordHandler(db *bun.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromContext(r.Context())
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, ErrorResponse{
				Status:  "error",
				Message: "Unauthorized",
			})
			return
		}

		var req ChangePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonResponse(w, http.StatusBadRequest, ErrorResponse{
				Status:  "error",
				Message: "Invalid request payload",
			})
			return
		}

		if len(req.NewPassword) < 8 {
			jsonResponse(w, http.StatusBadRequest, ErrorResponse{
				Status:  "error",
				Message: "New password must be at least 8 characters long",
			})
			return
		}

		var staff models.Staff
		err = db.NewSelect().Model(&staff).Where("id = ?", userID).Scan(r.Context())
		if err != nil {
			jsonResponse(w, http.StatusNotFound, ErrorResponse{
				Status:  "error",
				Message: "User not found",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(req.OldPassword)); err != nil {
			jsonResponse(w, http.StatusUnauthorized, ErrorResponse{
				Status:  "error",
				Message: "Old password is incorrect",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Failed to hash new password",
			})
			return
		}

		_, err = db.NewUpdate().
			Model(&staff).
			Set("password = ?", string(hashedPassword)).
			Where("id = ?", userID).
			Exec(r.Context())

		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Failed to update password",
			})
			return
		}

		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Password updated successfully",
		})
	}
}
