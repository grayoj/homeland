package routes

import (
	"homeland/handlers/staff"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func RegisterStaffRoutes(r chi.Router, db *bun.DB) {
	r.Route("/staff", func(r chi.Router) {
		r.Get("/{id}", staff.GetStaffHandler(db))
		r.Put("/{id}", staff.UpdateStaffHandler(db))
		r.Delete("/{id}", staff.DeleteStaffHandler(db))
	})
}
