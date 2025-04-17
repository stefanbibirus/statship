package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	
	"relationship-helix/internal/api/handlers"
	"relationship-helix/internal/api/middleware"
	"relationship-helix/internal/config"
)

// SetupRoutes configurează rutele API
func SetupRoutes(app *fiber.App, db *sql.DB, cfg *config.Config) {
	// Creează handler-ele
	authHandler := handlers.NewAuthHandler(db, cfg)
	relationshipHandler := handlers.NewRelationshipHandler(db, cfg)
	
	// Grupul de rute API
	api := app.Group("/api")
	
	// Rute de autentificare (publice)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	
	// Rute protejate prin autentificare
	auth.Get("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.GetMe)
	
	// Rute pentru relații (protejate)
	relationship := api.Group("/relationship", middleware.AuthMiddleware(cfg.JWTSecret))
	relationship.Get("/", relationshipHandler.GetRelationship)
	relationship.Post("/invite", relationshipHandler.GenerateInviteCode)
	relationship.Post("/join", relationshipHandler.UseInviteCode)
	relationship.Post("/position", relationshipHandler.UpdatePosition)
	relationship.Delete("/", relationshipHandler.DeleteRelationship)
}