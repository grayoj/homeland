package incident

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/uptrace/bun"
)

func GetIncidents(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var incidents []models.Incident

		err := db.NewSelect().Model(&incidents).Scan(ctx)
		if err != nil {
			log.Printf("DB error: %v", err)

			if ctx.Err() == context.DeadlineExceeded {
				utils.RespondWithError(w, http.StatusGatewayTimeout, "Database request timed out")
				return
			}

			if errors.Is(err, sql.ErrConnDone) {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Database connection lost")
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithJSON(w, http.StatusOK, []models.Incident{})
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "Unexpected database error")
			return
		}

		if len(incidents) == 0 {
			utils.RespondWithJSON(w, http.StatusOK, []models.Incident{})
			return
		}

		// Return the retrieved incidents
		utils.RespondWithJSON(w, http.StatusOK, incidents)
	}
}
