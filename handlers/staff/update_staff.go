package staff

import (
	"encoding/json"
	"net/http"
	"time"

	"homeland/models"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

var allowedRoles = map[models.RoleEnum]bool{
	models.RoleAdmin:    true,
	models.RoleSSA:      true,
	models.RoleDirector: true,
}

func isAuthorized(role models.RoleEnum) bool {
	return allowedRoles[role]
}

func UpdateStaffHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		staffID := chi.URLParam(r, "id")
		var req OnboardRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		var staff models.Staff
		err := db.NewSelect().Model(&staff).Where("id = ?", staffID).Scan(r.Context())
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Staff not found")
			return
		}

		if !isAuthorized(models.RoleEnum(req.Role)) {
			respondWithError(w, http.StatusForbidden, "You are not authorized to update staff")
			return
		}

		staff.FirstName = req.FirstName
		staff.MiddleName = req.MiddleName
		staff.LastName = req.LastName
		staff.Position = models.PositionEnum(req.Position)
		staff.Address = req.Address
		staff.Department = models.DepartmentEnum(req.Department)
		staff.StateOfOrigin = req.StateOfOrigin
		staff.Role = models.RoleEnum(req.Role)
		staff.UpdatedAt = time.Now()

		_, err = db.NewUpdate().Model(&staff).Where("id = ?", staffID).Exec(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update staff")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Staff updated successfully",
			"data":    staff,
		})
	}
}
