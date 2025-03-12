package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	routes "homeland/api"
	"homeland/config"
	"homeland/middleware"
	"homeland/models"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	seedAdmin(db, cfg)

	r := chi.NewRouter()

	r.Use(middleware.Logging)

	r.Route("/api/v1", func(r chi.Router) {
		routes.RegisterAuthRoutes(r, db, cfg)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTMiddleware(cfg.JWTSecret))

			routes.RegisterAdminRoutes(r, db, cfg)
			routes.RegisterStaffRoutes(r, db)
		})
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func createTables(db *bun.DB) error {
	ctx := context.Background()

	_, err := db.NewCreateTable().
		Model((*models.Staff)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("Tables created or already exist.")
	return nil
}

func seedAdmin(db *bun.DB, cfg *config.Config) {
	ctx := context.Background()

	count, err := db.NewSelect().
		Model((*models.Staff)(nil)).
		Where("email = ?", cfg.AdminEmail).
		Count(ctx)
	if err != nil {
		log.Printf("Error checking for admin existence: %v", err)
		return
	}

	if count == 0 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}

		admin := models.Staff{
			FirstName:     "Admin",
			MiddleName:    "",
			LastName:      "User",
			Email:         cfg.AdminEmail,
			Password:      string(hashed),
			AgentID:       "ADMIN001",
			ProfilePhoto:  "",
			Position:      models.PositionIT,
			Address:       "Head Office",
			Department:    models.DeptHomelandSecurity,
			DateOfBirth:   time.Now(),
			StateOfOrigin: "Abia",
			Role:          models.RoleAdmin,
		}

		_, err = db.NewInsert().Model(&admin).Exec(ctx)
		if err != nil {
			log.Printf("Error seeding admin: %v", err)
			return
		}

		log.Printf("Admin account seeded with email: %s", cfg.AdminEmail)
	} else {
		log.Println("Admin account already exists; skipping seeding")
	}
}
