package handlers

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"

	"relationship-helix/internal/config"
	"relationship-helix/internal/models"
	"relationship-helix/internal/utils"
)

// RelationshipHandler gestionează rutele de relații
type RelationshipHandler struct {
	DB     *sql.DB
	Config *config.Config
}

// NewRelationshipHandler creează un nou handler de relații
func NewRelationshipHandler(db *sql.DB, cfg *config.Config) *RelationshipHandler {
	return &RelationshipHandler{
		DB:     db,
		Config: cfg,
	}
}

// GetRelationship returnează relația utilizatorului curent
func (h *RelationshipHandler) GetRelationship(c *fiber.Ctx) error {
	// Obține ID-ul utilizatorului din context
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Neautentificat",
		})
	}
	
	// Caută relația utilizatorului
	var relationship models.Relationship
	err := h.DB.QueryRow(
		`SELECT id, user1_id, user2_id, user1_name, user2_name, start_date, created_at, updated_at 
         FROM relationships 
         WHERE user1_id = $1 OR user2_id = $1`,
		userID,
	).Scan(
		&relationship.ID, 
		&relationship.User1ID, 
		&relationship.User2ID, 
		&relationship.User1Name, 
		&relationship.User2Name, 
		&relationship.StartDate, 
		&relationship.CreatedAt, 
		&relationship.UpdatedAt,
	)
	
	// Dacă nu există nicio relație
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"relationship": nil,
		})
	}
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea relației",
		})
	}
	
	// Obține pozițiile curbelor
	var userPosition, partnerPosition int
	
	// Poziția utilizatorului
	err = h.DB.QueryRow(
		`SELECT position 
         FROM curve_positions 
         WHERE relationship_id = $1 AND user_id = $2`,
		relationship.ID, userID,
	).Scan(&userPosition)
	
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea poziției utilizatorului",
		})
	}
	
	// Poziția partenerului
	var partnerID uint
	if relationship.User1ID == userID {
		partnerID = relationship.User2ID
	} else {
		partnerID = relationship.User1ID
	}
	
	err = h.DB.QueryRow(
		`SELECT position 
         FROM curve_positions 
         WHERE relationship_id = $1 AND user_id = $2`,
		relationship.ID, partnerID,
	).Scan(&partnerPosition)
	
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea poziției partenerului",
		})
	}
	
	// Returnează relația și pozițiile
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"relationship":         relationship.ToResponse(userID),
		"userCurvePosition":    userPosition,
		"partnerCurvePosition": partnerPosition,
	})
}

// GenerateInviteCode generează un cod de invitație
func (h *RelationshipHandler) GenerateInviteCode(c *fiber.Ctx) error {
	// Obține ID-ul utilizatorului din context
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Neautentificat",
		})
	}
	
	// Verifică dacă utilizatorul are deja o relație
	var exists bool
	err := h.DB.QueryRow(
		`SELECT EXISTS(
            SELECT 1 FROM relationships 
            WHERE user1_id = $1 OR user2_id = $1
         )`,
		userID,
	).Scan(&exists)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la verificarea relațiilor existente",
		})
	}
	
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Ai deja o relație activă",
		})
	}
	
	// Generează codul de invitație
	code, err := utils.GenerateInviteCode()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la generarea codului de invitație",
		})
	}
	
	// Calculează data de expirare
	expiresAt := time.Now().Add(h.Config.InviteCodeExpiration)
	
	// Șterge orice cod de invitație existent
	_, err = h.DB.Exec(
		`DELETE FROM invite_codes WHERE user_id = $1`,
		userID,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la ștergerea codurilor existente",
		})
	}
	
	// Salvează codul de invitație
	_, err = h.DB.Exec(
		`INSERT INTO invite_codes (user_id, code, expires_at, created_at) 
         VALUES ($1, $2, $3, NOW())`,
		userID, code, expiresAt,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la salvarea codului de invitație",
		})
	}
	
	// Returnează codul de invitație
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"inviteCode": code,
		"expiresAt":  expiresAt,
	})
}

// UseInviteCodeRequest reprezintă cererea de utilizare a unui cod de invitație
type UseInviteCodeRequest struct {
	InviteCode string `json:"inviteCode" validate:"required"`
}

// UseInviteCode utilizează un cod de invitație pentru a crea o relație
func (h *RelationshipHandler) UseInviteCode(c *fiber.Ctx) error {
	// Obține ID-ul utilizatorului din context
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Neautentificat",
		})
	}
	
	// Parsează cererea
	var req UseInviteCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cerere invalidă",
		})
	}
	
	// Verifică codul de invitație
	if req.InviteCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Codul de invitație este obligatoriu",
		})
	}
	
	// Verifică dacă utilizatorul are deja o relație
	var hasRelationship bool
	err := h.DB.QueryRow(
		`SELECT EXISTS(
            SELECT 1 FROM relationships 
            WHERE user1_id = $1 OR user2_id = $1
         )`,
		userID,
	).Scan(&hasRelationship)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la verificarea relațiilor existente",
		})
	}
	
	if hasRelationship {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Ai deja o relație activă",
		})
	}
	
	// Începe o tranzacție
	tx, err := h.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la inițierea tranzacției",
		})
	}
	defer tx.Rollback()
	
	// Caută codul de invitație
	var inviteCode models.InviteCode
	err = tx.QueryRow(
		`SELECT id, user_id, code, expires_at, created_at 
         FROM invite_codes 
         WHERE code = $1 AND expires_at > NOW()`,
		req.InviteCode,
	).Scan(&inviteCode.ID, &inviteCode.UserID, &inviteCode.Code, &inviteCode.ExpiresAt, &inviteCode.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Cod de invitație invalid sau expirat",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la căutarea codului de invitație",
		})
	}
	
	// Verifică dacă codul nu aparține utilizatorului curent
	if inviteCode.UserID == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Nu poți folosi propriul cod de invitație",
		})
	}
	
	// Verifică dacă partenerul are deja o relație
	var partnerHasRelationship bool
	err = tx.QueryRow(
		`SELECT EXISTS(
            SELECT 1 FROM relationships 
            WHERE user1_id = $1 OR user2_id = $1
         )`,
		inviteCode.UserID,
	).Scan(&partnerHasRelationship)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la verificarea relațiilor partenerului",
		})
	}
	
	if partnerHasRelationship {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Partenerul are deja o relație activă",
		})
	}
	
	// Obține informații despre utilizatorul curent și partener
	var currentUsername, partnerUsername string
	
	err = tx.QueryRow(
		`SELECT username FROM users WHERE id = $1`,
		userID,
	).Scan(&currentUsername)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea informațiilor utilizatorului",
		})
	}
	
	err = tx.QueryRow(
		`SELECT username FROM users WHERE id = $1`,
		inviteCode.UserID,
	).Scan(&partnerUsername)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea informațiilor partenerului",
		})
	}
	
	// Creează relația
	var relationshipID uint
	err = tx.QueryRow(
		`INSERT INTO relationships (user1_id, user2_id, user1_name, user2_name, start_date, created_at, updated_at) 
         VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW()) 
         RETURNING id`,
		inviteCode.UserID, userID, partnerUsername, currentUsername,
	).Scan(&relationshipID)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la crearea relației",
		})
	}
	
	// Inițializează pozițiile curbelor
	_, err = tx.Exec(
		`INSERT INTO curve_positions (relationship_id, user_id, position, created_at, updated_at) 
         VALUES ($1, $2, 0, NOW(), NOW()), ($1, $3, 0, NOW(), NOW())`,
		relationshipID, userID, inviteCode.UserID,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la inițializarea pozițiilor curbelor",
		})
	}
	
	// Șterge codul de invitație
	_, err = tx.Exec(
		`DELETE FROM invite_codes WHERE id = $1`,
		inviteCode.ID,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la ștergerea codului de invitație",
		})
	}
	
	// Commit tranzacția
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la finalizarea tranzacției",
		})
	}
	
	// Obține relația creată
	var relationship models.Relationship
	err = h.DB.QueryRow(
		`SELECT id, user1_id, user2_id, user1_name, user2_name, start_date, created_at, updated_at 
         FROM relationships 
         WHERE id = $1`,
		relationshipID,
	).Scan(
		&relationship.ID, 
		&relationship.User1ID, 
		&relationship.User2ID, 
		&relationship.User1Name, 
		&relationship.User2Name, 
		&relationship.StartDate, 
		&relationship.CreatedAt, 
		&relationship.UpdatedAt,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea relației create",
		})
	}
	
	// Returnează relația creată
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"relationship":         relationship.ToResponse(userID),
		"userCurvePosition":    0,
		"partnerCurvePosition": 0,
	})
}

// UpdatePositionRequest reprezintă cererea de actualizare a poziției
type UpdatePositionRequest struct {
	Position int `json:"position" validate:"required,min=0,max=100"`
}

// UpdatePosition actualizează poziția curbei utilizatorului
func (h *RelationshipHandler) UpdatePosition(c *fiber.Ctx) error {
	// Obține ID-ul utilizatorului din context
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Neautentificat",
		})
	}
	
	// Parsează cererea
	var req UpdatePositionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Cerere invalidă",
		})
	}
	
	// Validează poziția
	if req.Position < 0 || req.Position > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Poziția trebuie să fie între 0 și 100",
		})
	}
	
	// Obține relația utilizatorului
	var relationship models.Relationship
	err := h.DB.QueryRow(
		`SELECT id, user1_id, user2_id 
         FROM relationships 
         WHERE user1_id = $1 OR user2_id = $1`,
		userID,
	).Scan(&relationship.ID, &relationship.User1ID, &relationship.User2ID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Nu ai o relație activă",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea relației",
		})
	}
	
	// Actualizează poziția curbei
	_, err = h.DB.Exec(
		`INSERT INTO curve_positions (relationship_id, user_id, position, created_at, updated_at) 
         VALUES ($1, $2, $3, NOW(), NOW()) 
         ON CONFLICT (relationship_id, user_id) 
         DO UPDATE SET position = $3, updated_at = NOW()`,
		relationship.ID, userID, req.Position,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la actualizarea poziției",
		})
	}
	
	// Determină ID-ul partenerului
	var partnerID uint
	if relationship.User1ID == userID {
		partnerID = relationship.User2ID
	} else {
		partnerID = relationship.User1ID
	}
	
	// Trimite notificare prin WebSocket
	update := models.PositionUpdate{
		RelationshipID: relationship.ID,
		UserID:         userID,
		PartnerID:      partnerID,
		Position:       req.Position,
	}
	
	// Trimite actualizarea poziției
	BroadcastPositionUpdate(update)
	
	// Returnează succes
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"position": req.Position,
	})
}

// DeleteRelationship șterge relația utilizatorului
func (h *RelationshipHandler) DeleteRelationship(c *fiber.Ctx) error {
	// Obține ID-ul utilizatorului din context
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Neautentificat",
		})
	}
	
	// Începe o tranzacție
	tx, err := h.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la inițierea tranzacției",
		})
	}
	defer tx.Rollback()
	
	// Obține relația utilizatorului
	var relationshipID uint
	err = tx.QueryRow(
		`SELECT id FROM relationships WHERE user1_id = $1 OR user2_id = $1`,
		userID,
	).Scan(&relationshipID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Nu ai o relație activă",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la obținerea relației",
		})
	}
	
	// Șterge pozițiile curbelor
	_, err = tx.Exec(
		`DELETE FROM curve_positions WHERE relationship_id = $1`,
		relationshipID,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la ștergerea pozițiilor curbelor",
		})
	}
	
	// Șterge relația
	_, err = tx.Exec(
		`DELETE FROM relationships WHERE id = $1`,
		relationshipID,
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la ștergerea relației",
		})
	}
	
	// Commit tranzacția
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Eroare la finalizarea tranzacției",
		})
	}
	
	// Returnează succes
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}