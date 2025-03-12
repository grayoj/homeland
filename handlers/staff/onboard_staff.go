package staff

import (
	"encoding/json"
	"net/http"
	"time"

	"homeland/config"
	"homeland/models"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type OnboardRequest struct {
	FirstName     string `json:"first_name"`
	MiddleName    string `json:"middle_name,omitempty"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	AgentID       string `json:"agent_id"`
	ProfilePhoto  string `json:"profile_photo,omitempty"`
	Position      string `json:"position"`
	Address       string `json:"address"`
	Department    string `json:"department"`
	DateOfBirth   string `json:"date_of_birth"`
	StateOfOrigin string `json:"state_of_origin"`
	Role          string `json:"role"`
}

func OnboardStaffHandler(db *bun.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OnboardRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		var existingStaff models.Staff
		err := db.NewSelect().Model(&existingStaff).Where("email = ?", req.Email).Scan(r.Context())
		if err == nil {
			respondWithError(w, http.StatusConflict, "User with this email already exists")
			return
		}

		if len(req.Password) < 8 {
			respondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
			return
		}

		parsedDOB, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		staff := models.Staff{
			FirstName:          req.FirstName,
			MiddleName:         req.MiddleName,
			LastName:           req.LastName,
			Email:              req.Email,
			Password:           string(hashedPassword),
			AgentID:            req.AgentID,
			ProfilePhoto:       req.ProfilePhoto,
			Position:           models.PositionEnum(req.Position),
			Address:            req.Address,
			Department:         models.DepartmentEnum(req.Department),
			DateOfBirth:        parsedDOB,
			StateOfOrigin:      req.StateOfOrigin,
			Role:               models.RoleEnum(req.Role),
			MustChangePassword: false,
		}

		_, err = db.NewInsert().Model(&staff).Exec(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to onboard staff")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Staff onboarded successfully",
			"data": map[string]interface{}{
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
