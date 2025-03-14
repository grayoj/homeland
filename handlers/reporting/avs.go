package reporting

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

func CreateAVSReport(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())

		if user.Department != string(models.DeptAVS) {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized: Only AVS can report incidents")
			return
		}

		var report models.AVSReport
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		report.ReportedBy = user.Email
		report.DateReported = time.Now()
		report.Department = string(models.DeptAVS)

		_, err := db.NewInsert().Model(&report).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create AVS report")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, report)
	}
}

func GetAVSReports(db *bun.DB) http.HandlerFunc {
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

		var reports []models.AVSReport

		err = db.NewSelect().Model(&reports).
			Limit(limit).
			Offset(offset).
			Order("date_reported DESC").
			Scan(ctx)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch AVS reports")
			return
		}

		total, err := db.NewSelect().Model((*models.AVSReport)(nil)).Count(ctx)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve report count")
			return
		}

		response := map[string]interface{}{
			"status":  "success",
			"message": "AVS reports retrieved successfully",
			"data":    reports,
			"pagination": map[string]int{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

func GetAVSReportByID(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var report models.AVSReport
		err := db.NewSelect().Model(&report).Where("id = ?", id).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "AVS report not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, report)
	}
}
