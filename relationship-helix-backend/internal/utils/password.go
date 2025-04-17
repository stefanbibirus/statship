package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generează un hash pentru o parolă folosind bcrypt
func HashPassword(password string) (string, error) {
	// Generează hash cu costul implicit (10)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(bytes), nil
}

// CheckPassword verifică dacă o parolă corespunde hash-ului său
func CheckPassword(hashedPassword, password string) error {
	// Compară parola cu hash-ul
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}