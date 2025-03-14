package workplace

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

var allowedRoles = map[models.RoleEnum]bool{
	models.RoleAdmin:    true,
	models.RoleSSA:      true,
	models.RoleDirector: true,
	models.RoleStaff:    true,
}

func CreateAppointment(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())
		if user.Department != string(models.DeptHomelandSecurity) || !allowedRoles[models.RoleEnum(user.Role)] {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized")
			return
		}

		var appointment models.Appointment
		if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		appointment.CreatedAt = time.Now()
		appointment.UpdatedAt = time.Now()

		_, err := db.NewInsert().Model(&appointment).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create appointment")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, appointment)
	}
}

func GetAppointments(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 {
			offset = 0
		}

		var appointments []models.Appointment

		err = db.NewSelect().Model(&appointments).
			Limit(limit).
			Offset(offset).
			Order("created_at DESC").
			Scan(ctx)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch appointments")
			return
		}

		total, err := db.NewSelect().Model((*models.Appointment)(nil)).Count(ctx)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch appointment count")
			return
		}

		response := map[string]interface{}{
			"data":       appointments,
			"pagination": map[string]int{"total": total, "limit": limit, "offset": offset},
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

func GetAppointmentByID(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
			return
		}

		var appointment models.Appointment
		err = db.NewSelect().Model(&appointment).Where("id = ?", id).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Appointment not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, appointment)
	}
}

func UpdateAppointment(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())
		if user.Department != string(models.DeptHomelandSecurity) || !allowedRoles[models.RoleEnum(user.Role)] {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized")
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
			return
		}

		var appointment models.Appointment
		err = db.NewSelect().Model(&appointment).Where("id = ?", id).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Appointment not found")
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		appointment.UpdatedAt = time.Now()
		_, err = db.NewUpdate().Model(&appointment).Where("id = ?", id).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update appointment")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, appointment)
	}
}

func DeleteAppointment(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())
		if user.Department != string(models.DeptHomelandSecurity) || !allowedRoles[models.RoleEnum(user.Role)] {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized")
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
			return
		}

		_, err = db.NewDelete().Model((*models.Appointment)(nil)).Where("id = ?", id).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete appointment")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Appointment deleted"})
	}
}
