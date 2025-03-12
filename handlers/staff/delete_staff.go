package staff

import (
	"encoding/json"
	"homeland/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func DeleteStaffHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		staffID := chi.URLParam(r, "id")

		_, err := db.NewDelete().Model((*models.Staff)(nil)).Where("id = ?", staffID).Exec(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete staff")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Staff deleted successfully",
		})
	}
}
