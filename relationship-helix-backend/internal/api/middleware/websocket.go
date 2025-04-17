package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	
	"relationship-helix/internal/utils"
)

// WebsocketAuth verifică autentificarea pentru conexiunile WebSocket
func WebsocketAuth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Verifică dacă cererea este pentru upgrade la WebSocket
		if websocket.IsWebSocketUpgrade(c) {
			// Obține token-ul din query params
			token := c.Query("token")
			
			// Verifică dacă token-ul există
			if token == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"message": "Token de autentificare lipsă",
				})
			}
			
			// Validează token-ul
			userID, err := utils.ValidateToken(token, jwtSecret)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"message": "Token invalid: " + err.Error(),
				})
			}
			
			// Setează ID-ul utilizatorului în locals pentru a fi utilizat în handler-ul WebSocket
			c.Locals("userID", userID)
			
			// Continuă cu upgrade-ul WebSocket
			return c.Next()
		}
		
		// Dacă nu este upgrade WebSocket, continuă normal
		return c.Next()
	}
}