package staff

import (
	"encoding/json"
	"log"
	"net/http"

	"homeland/models"
	"homeland/utils"

	"github.com/uptrace/bun"
)

func GetAllStaffHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var staffList []models.Staff

		err := db.NewSelect().Model(&staffList).Scan(r.Context())
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve staff records")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Staff records retrieved successfully",
			"data":    staffList,
		})
	}
}
