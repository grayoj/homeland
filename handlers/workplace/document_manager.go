package workplace

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"homeland/models"
	"homeland/utils"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

var allowedUploadPositions = map[models.PositionEnum]bool{
	models.PositionSSA:        true,
	models.PositionIT:         true,
	models.PositionHR:         true,
	models.PositionDirector:   true,
	models.PositionCallCenter: true,
}

func UploadDocument(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r.Context())

		if user.Department != string(models.DeptHomelandSecurity) {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized: Only Homeland Security can upload documents")
			return
		}

		var doc models.Document
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		doc.UploadedBy = user.Email
		doc.CreatedAt = time.Now()

		_, err := db.NewInsert().Model(&doc).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to upload document")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, doc)
	}
}

func GetDocuments(db *bun.DB) http.HandlerFunc {
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

		var docs []models.Document

		err = db.NewSelect().Model(&docs).
			Limit(limit).
			Offset(offset).
			Order("created_at DESC").
			Scan(ctx)

		if err != nil {
			log.Printf("DB error: %v", err)

			if ctx.Err() == context.DeadlineExceeded {
				utils.RespondWithError(w, http.StatusGatewayTimeout, "Database request timed out. Please try again later.")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "An unexpected error occurred while fetching documents.")
			return
		}

		total, err := db.NewSelect().Model((*models.Document)(nil)).Count(ctx)
		if err != nil {
			log.Printf("Count query error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch document count.")
			return
		}

		response := map[string]interface{}{
			"data":       docs,
			"pagination": map[string]int{"total": total, "limit": limit, "offset": offset},
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

func GetDocumentByID(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var doc models.Document
		err := db.NewSelect().Model(&doc).Where("id = ?", id).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Document not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, doc)
	}
}
