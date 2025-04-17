package models

import "time"

// CurvePosition reprezintă poziția curbei unui utilizator în animația helixului
type CurvePosition struct {
	ID             uint      `json:"id"`
	RelationshipID uint      `json:"relationshipId"`
	UserID         uint      `json:"userId"`
	Position       int       `json:"position"` // 0-100, unde 0 = apropiat, 100 = distant
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// PositionUpdate reprezintă o actualizare de poziție trimisă prin WebSocket
type PositionUpdate struct {
	RelationshipID uint `json:"relationshipId"`
	UserID         uint `json:"userId"`
	PartnerID      uint `json:"partnerId"`
	Position       int  `json:"position"`
}