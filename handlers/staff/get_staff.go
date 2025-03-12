package staff

import (
	"encoding/json"
	"net/http"

	"homeland/models"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func GetStaffHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		staffID := chi.URLParam(r, "id")

		var staff models.Staff
		err := db.NewSelect().Model(&staff).Where("id = ?", staffID).Scan(r.Context())
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Staff not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Staff retrieved successfully",
			"data":    staff,
		})
	}
}
