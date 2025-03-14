package staff

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/uptrace/bun"
)

func GetAllStaffHandler(db *bun.DB) http.HandlerFunc {
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

		var staffList []models.Staff

		err = db.NewSelect().Model(&staffList).
			Limit(limit).
			Offset(offset).
			Order("created_at DESC").
			Scan(ctx)

		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve staff records")
			return
		}

		total, err := db.NewSelect().Model((*models.Staff)(nil)).Count(ctx)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve staff count")
			return
		}

		response := map[string]interface{}{
			"status":  "success",
			"message": "Staff records retrieved successfully",
			"data":    staffList,
			"pagination": map[string]int{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}
