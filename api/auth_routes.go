package routes

import (
	"homeland/config"
	"homeland/handlers/auth"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func RegisterAuthRoutes(r chi.Router, db *bun.DB, cfg *config.Config) {
	r.Post("/login", auth.LoginHandler(db, cfg))
	r.Post("/refresh", auth.RefreshTokenHandler(cfg))
	r.Get("/auth", auth.AuthCheckHandler(db, cfg))
}
