package api

import (
	"homeland/handlers/incident"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func RegisterIncidentRoutes(r chi.Router, db *bun.DB) {
	r.Post("/incidents", incident.CreateIncidentHandler(db))
	r.Get("/incidents", incident.GetIncidents(db))
	r.Get("/incidents/{id}", incident.GetIncidentByID(db))
	r.Put("/incidents/{id}", incident.UpdateIncident(db))
	r.Delete("/incidents/{id}", incident.DeleteIncident(db))
}
