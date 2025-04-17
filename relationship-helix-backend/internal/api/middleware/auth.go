package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	
	"relationship-helix/internal/utils"
)

// AuthMiddleware verifică și validează token-ul JWT din header-ul Authorization
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obține header-ul Authorization
		authHeader := c.Get("Authorization")
		
		// Verifică dacă header-ul există
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token de autentificare lipsă",
			})
		}
		
		// Extrage token-ul din header
		// Format: "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Format token invalid",
			})
		}
		
		tokenString := parts[1]
		
		// Validează token-ul
		userID, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token invalid: " + err.Error(),
			})
		}
		
		// Setează ID-ul utilizatorului în context pentru a fi utilizat în handler-e
		c.Locals("userID", userID)
		
		// Continuă cu cererea
		return c.Next()
	}
}