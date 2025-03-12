package incident

import (
	"context"
	"log"
	"net/http"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func DeleteIncident(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		id := chi.URLParam(r, "id")
		if id == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Incident ID is required")
			return
		}

		res, err := db.NewDelete().
			Model((*models.Incident)(nil)).
			Where("id = ?", id).
			Exec(ctx)

		if err != nil {
			log.Printf("DB error: %v", err)

			if ctx.Err() == context.DeadlineExceeded {
				utils.RespondWithError(w, http.StatusGatewayTimeout, "Database request timed out")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete incident")
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			utils.RespondWithError(w, http.StatusNotFound, "Incident not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Incident deleted successfully"})
	}
}
