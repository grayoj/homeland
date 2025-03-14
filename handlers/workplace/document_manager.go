package workplace

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
		var docs []models.Document
		err := db.NewSelect().Model(&docs).Scan(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch documents")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, docs)
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
