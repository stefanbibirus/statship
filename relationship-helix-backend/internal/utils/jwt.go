package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken generează un token JWT pentru autentificare
func GenerateToken(userID uint, secret string, expiration time.Duration) (string, error) {
	// Crează un token nou cu algoritmul de semnare HS256
	token := jwt.New(jwt.SigningMethodHS256)
	
	// Setează claims (revendicări)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["exp"] = time.Now().Add(expiration).Unix()
	
	// Semnează tokenul cu cheia secretă
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken verifică dacă un token JWT este valid
func ValidateToken(tokenString string, secret string) (uint, error) {
	// Parsează tokenul
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifică metoda de semnare
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodă de semnare invalidă")
		}
		
		// Returnează cheia secretă pentru verificare
		return []byte(secret), nil
	})
	
	// Verifică erorile de parsare
	if err != nil {
		return 0, err
	}
	
	// Verifică dacă tokenul este valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verifică claims
		if userID, ok := claims["id"].(float64); ok {
			return uint(userID), nil
		}
		
		return 0, errors.New("ID utilizator invalid în token")
	}
	
	return 0, errors.New("token invalid")
}