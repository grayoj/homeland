package api

import (
	"homeland/handlers/reporting"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func RegisterReportingRoutes(r chi.Router, db *bun.DB) {
	r.Route("/reports", func(r chi.Router) {
		r.Route("/fire", func(r chi.Router) {
			r.Post("/", reporting.CreateFireReport(db))
			r.Get("/", reporting.GetFireReports(db))
			r.Get("/{id}", reporting.GetFireReportByID(db))
		})

		r.Route("/ems", func(r chi.Router) {
			r.Post("/", reporting.CreateEMSReport(db))
			r.Get("/", reporting.GetEMSReports(db))
			r.Get("/{id}", reporting.GetEMSReportByID(db))
		})

		r.Route("/avs", func(r chi.Router) {
			r.Post("/", reporting.CreateAVSReport(db))
			r.Get("/", reporting.GetAVSReports(db))
			r.Get("/{id}", reporting.GetAVSReportByID(db))
		})
	})
}
