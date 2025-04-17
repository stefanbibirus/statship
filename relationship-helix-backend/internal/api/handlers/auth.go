package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"relationship-helix/internal/config"
	"relationship-helix/internal/models"
	"relationship-helix/internal/utils"
)

// AuthHandler gestionează rutele de autentificare
type AuthHandler struct {
	DB     *sql.DB
	Config *config.Config
}

// NewAuthHandler creează un nou handler de autentificare
func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		DB:     db,
		Config: cfg,
	}
}

// RegisterRequest reprezintă cererea de înregistrare
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Register înregistrează un nou utilizator
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parsează cererea
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cerere invalidă",
		})
	}
	
	// Validează câmpurile
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Toate câmpurile sunt obligatorii",
		})
	}
	
	// Verifică dacă email-ul există deja
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la verificarea email-ului",
		})
	}
	
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Email-ul este deja utilizat",
		})
	}
	
	// Hash-uiește parola
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la hash-uirea parolei",
		})
	}
	
	// Creează utilizatorul
	var user models.User
	err = h.DB.QueryRow(
		`INSERT INTO users (username, email, password, created_at, updated_at) 
         VALUES ($1, $2, $3, NOW(), NOW()) 
         RETURNING id, username, email, created_at, updated_at`,
		req.Username, req.Email, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la crearea utilizatorului",
		})
	}
	
	// Generează token JWT
	token, err := utils.GenerateToken(user.ID, h.Config.JWTSecret, h.Config.JWTExpiration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la generarea token-ului",
		})
	}
	
	// Returnează utilizatorul și token-ul
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// LoginRequest reprezintă cererea de autentificare
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login autentifică un utilizator existent
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parsează cererea
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cerere invalidă",
		})
	}
	
	// Validează câmpurile
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Email-ul și parola sunt obligatorii",
		})
	}
	
	// Caută utilizatorul după email
	var user models.User
	err := h.DB.QueryRow(
		`SELECT id, username, email, password, created_at, updated_at 
         FROM users 
         WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Email sau parolă invalidă",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la căutarea utilizatorului",
		})
	}
	
	// Verifică parola
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Email sau parolă invalidă",
		})
	}
	
	// Generează token JWT
	token, err := utils.GenerateToken(user.ID, h.Config.JWTSecret, h.Config.JWTExpiration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la generarea token-ului",
		})
	}
	
	// Returnează utilizatorul și token-ul
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// GetMe returnează informațiile utilizatorului curent
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Obține ID-ul utilizatorului din context (setat de middleware-ul de autentificare)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Neautentificat",
		})
	}
	
	// Caută utilizatorul în baza de date
	var user models.User
	err := h.DB.QueryRow(
		`SELECT id, username, email, created_at, updated_at 
         FROM users 
         WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Utilizatorul nu a fost găsit",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea informațiilor utilizatorului",
		})
	}
	
	// Returnează utilizatorul
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user.ToResponse(),
	})
}