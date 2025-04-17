package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"strconv"

	"github.com/gofiber/websocket/v2"

	"relationship-helix/internal/models"
)

// Map pentru a ține evidența conexiunilor WebSocket
// Map[relationshipID]Map[userID]*websocket.Conn
var (
	clientsMutex sync.RWMutex
	clients      = make(map[uint]map[uint]*websocket.Conn)
)

// HandleWebsocketConnection gestionează o conexiune WebSocket
func HandleWebsocketConnection(c *websocket.Conn) {
	// Obține ID-ul utilizatorului și ID-ul relației din URL
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		log.Println("WebSocket: ID utilizator invalid")
		return
	}
	
	// Parametri din URL (/ws/relationship/:id)
	relationshipID := c.Params("id") // Obținem relationshipID direct din parametrul URL
	
	if relationshipID == "" {
		log.Println("WebSocket: ID relație lipsă")
		return
	}
	
	// Convertim relationshipID din string în uint
	relID, err := strconv.ParseUint(relationshipID, 10, 64)
	if err != nil {
		log.Printf("WebSocket: ID relație invalid: %v\n", err)
		return
	}
	
	// Acum avem relID ca uint64, îl convertim la uint
	relIDUint := uint(relID)
	
	// Adaugă clientul la hartă
	clientsMutex.Lock()
	if _, ok := clients[relIDUint]; !ok {
		clients[relIDUint] = make(map[uint]*websocket.Conn)
	}
	clients[relIDUint][userID] = c
	clientsMutex.Unlock()
	
	// Mesaj de conectare
	log.Printf("WebSocket: Utilizatorul %d s-a conectat la relația %d\n", userID, relIDUint)
	
	// Buclă de citire mesaje (nu este necesară pentru acest caz)
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket: Eroare la citirea mesajului: %v\n", err)
			break
		}
	}
	
	// Eliminare client la deconectare
	clientsMutex.Lock()
	delete(clients[relIDUint], userID)
	// Dacă nu mai există clienți pentru această relație, șterge și intrarea
	if len(clients[relIDUint]) == 0 {
		delete(clients, relIDUint)
	}
	clientsMutex.Unlock()
	
	log.Printf("WebSocket: Utilizatorul %d s-a deconectat de la relația %d\n", userID, relIDUint)
}

// BroadcastPositionUpdate trimite actualizări de poziție prin WebSocket
func BroadcastPositionUpdate(update models.PositionUpdate) {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	
	// Verifică dacă există conexiuni pentru această relație
	relationshipClients, ok := clients[update.RelationshipID]
	if !ok {
		return
	}
	
	// Construiește mesajul JSON
	message := map[string]interface{}{
		"type": "position_update",
		"payload": map[string]interface{}{
			"partnerId": update.UserID,
			"position":  update.Position,
		},
	}
	
	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("WebSocket: Eroare la serializarea mesajului: %v\n", err)
		return
	}
	
	// Trimite mesajul doar partenerului
	if conn, ok := relationshipClients[update.PartnerID]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("WebSocket: Eroare la trimiterea mesajului: %v\n", err)
		}
	}
}