package api

import (
	"homeland/handlers/workplace"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func RegisterWorkplaceRoutes(r chi.Router, db *bun.DB) {
	r.Route("/appointments", func(r chi.Router) {
		r.Post("/", workplace.CreateAppointment(db))
		r.Get("/", workplace.GetAppointments(db))
		r.Get("/{id}", workplace.GetAppointmentByID(db))
		r.Put("/{id}", workplace.UpdateAppointment(db))
		r.Delete("/{id}", workplace.DeleteAppointment(db))
	})
	r.Route("/documents", func(r chi.Router) {
		r.Post("/upload", workplace.UploadDocument(db))
		r.Get("/", workplace.GetDocuments(db))
		r.Get("/{id}", workplace.GetDocumentByID(db))
	})
}
