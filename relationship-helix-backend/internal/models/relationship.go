package models

import (
	"time"
)

// Relationship reprezintă o relație între doi utilizatori
type Relationship struct {
	ID              uint      `json:"id"`
	User1ID         uint      `json:"user1Id"`
	User2ID         uint      `json:"user2Id"`
	User1Name       string    `json:"user1Name"`
	User2Name       string    `json:"user2Name"`
	StartDate       time.Time `json:"startDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// RelationshipResponse este structura returnată în API
type RelationshipResponse struct {
	ID              uint      `json:"id"`
	PartnerID       uint      `json:"partnerId"`
	PartnerName     string    `json:"partnerName"`
	StartDate       time.Time `json:"startDate"`
	DaysSinceStart  int       `json:"daysSinceStart"`
}

// ToResponse convertește un Relationship într-un RelationshipResponse pentru utilizatorul specificat
func (r *Relationship) ToResponse(userID uint) RelationshipResponse {
	var partnerID uint
	var partnerName string
	
	if r.User1ID == userID {
		partnerID = r.User2ID
		partnerName = r.User2Name
	} else {
		partnerID = r.User1ID
		partnerName = r.User1Name
	}
	
	daysSinceStart := int(time.Since(r.StartDate).Hours() / 24)
	
	return RelationshipResponse{
		ID:              r.ID,
		PartnerID:       partnerID,
		PartnerName:     partnerName,
		StartDate:       r.StartDate,
		DaysSinceStart:  daysSinceStart,
	}
}

// InviteCode reprezintă un cod de invitație pentru a forma o relație
type InviteCode struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"userId"`
	Code         string    `json:"code"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}