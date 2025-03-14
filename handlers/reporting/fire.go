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

func CreateFireReport(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())

		if user.Department != string(models.DeptFireService) {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized: Only Fire Service can report incidents")
			return
		}

		var report models.FireReport
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		report.ReportedBy = user.Email
		report.DateReported = time.Now()
		report.Department = string(models.DeptFireService)

		_, err := db.NewInsert().Model(&report).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create fire report")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, report)
	}
}

func GetFireReports(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reports []models.FireReport
		err := db.NewSelect().Model(&reports).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch fire reports")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, reports)
	}
}

func GetFireReportByID(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var report models.FireReport
		err := db.NewSelect().Model(&report).Where("id = ?", id).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Fire report not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, report)
	}
}
