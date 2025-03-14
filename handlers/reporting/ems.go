package reporting

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func CreateEMSReport(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())

		if user.Department != string(models.DeptEMS) {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized: Only EMS can report incidents")
			return
		}

		var report models.EMSReport
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		report.ReportedBy = user.Email
		report.DateReported = time.Now()
		report.Department = string(models.DeptEMS)

		_, err := db.NewInsert().Model(&report).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create EMS report")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, report)
	}
}

func GetEMSReports(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reports []models.EMSReport
		err := db.NewSelect().Model(&reports).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch EMS reports")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, reports)
	}
}

func GetEMSReportByID(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var report models.EMSReport
		err := db.NewSelect().Model(&report).Where("id = ?", id).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "EMS report not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, report)
	}
}
