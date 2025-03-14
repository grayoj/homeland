package incident

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func GetIncidents(db *bun.DB) http.HandlerFunc {
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

		var incidents []models.Incident

		err = db.NewSelect().Model(&incidents).
			Limit(limit).
			Offset(offset).
			Order("date_reported DESC").
			Scan(ctx)

		if err != nil {
			log.Printf("DB error: %v", err)

			if ctx.Err() == context.DeadlineExceeded {
				utils.RespondWithError(w, http.StatusGatewayTimeout, "Database request timed out. Please try again later.")
				return
			}

			if errors.Is(err, sql.ErrConnDone) {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Database connection lost. Please refresh and try again.")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "An unexpected error occurred while fetching incidents.")
			return
		}

		total, err := db.NewSelect().Model((*models.Incident)(nil)).Count(ctx)
		if err != nil {
			log.Printf("Count query error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch incident count.")
			return
		}

		response := map[string]interface{}{
			"data":       incidents,
			"pagination": map[string]int{"total": total, "limit": limit, "offset": offset},
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

func GetIncidentByID(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := chi.URLParam(r, "id")
		if id == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Incident ID is required.")
			return
		}

		var incident models.Incident
		err := db.NewSelect().Model(&incident).Where("id = ?", id).Limit(1).Scan(ctx)

		if err != nil {
			log.Printf("DB error: %v", err)

			if ctx.Err() == context.DeadlineExceeded {
				utils.RespondWithError(w, http.StatusGatewayTimeout, "Database request timed out. Please try again later.")
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithError(w, http.StatusNotFound, "Incident not found. Please check the ID and try again.")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "An unexpected error occurred while fetching the incident.")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, incident)
	}
}
