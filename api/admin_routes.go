package api

import (
	"homeland/config"
	"homeland/handlers/auth"
	"homeland/handlers/staff"
	"homeland/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func RegisterAdminRoutes(r chi.Router, db *bun.DB, cfg *config.Config) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RoleAuthorization("Admin", "SSA", "Director"))
		r.Post("/onboard", staff.OnboardStaffHandler(db, cfg))
		r.Post("/change-password", auth.ChangePasswordHandler(db, cfg))
	})
}
