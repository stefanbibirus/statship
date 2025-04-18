package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"relationship-helix/internal/api/handlers"
	"relationship-helix/internal/api/middleware"
	"relationship-helix/internal/api/routes"
	"relationship-helix/internal/config"
	"relationship-helix/internal/db"
)

func main() {
	// Încarcă variabilele de mediu
	if err := godotenv.Load(); err != nil {
		log.Println("Nu s-a găsit fișierul .env, se folosesc variabilele de mediu existente")
	}

	// Inițializează configurația
	cfg := config.LoadConfig()

	// Inițializează conexiunea la baza de date
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Nu s-a putut conecta la baza de date: %v", err)
	}
	defer database.Close()

	// Rulează migrările
	//if err := db.RunMigrations(database); err != nil {
	//	log.Fatalf("Eroare la rularea migrărilor: %v", err)
	//}

	// Creează aplicația Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Gestionare personalizată a erorilor
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return ctx.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Middleware pentru WebSocket
	app.Use("/ws", middleware.WebsocketAuth(cfg.JWTSecret))

	// Setează rutele WebSocket
	app.Use("/ws/*", websocket.New(func(c *websocket.Conn) {
		// Conexiune WebSocket stabilită
		handlers.HandleWebsocketConnection(c)
	}))

	// Setează rutele API
	routes.SetupRoutes(app, database, cfg)

	// Determină portul serverului
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(cfg.ServerPort)
	}

	// Pornește serverul
	log.Printf("Serverul rulează pe portul %s\n", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Eroare la pornirea serverului: %v", err)
	}
}
