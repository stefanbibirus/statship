package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// Caractere pentru codul de invitație
	// Exclude caracterele similare (0/O, 1/I/l, etc.) pentru a preveni confuziile
	codeChars    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	codeLength   = 8
)

// GenerateInviteCode generează un cod aleatoriu pentru invitații
func GenerateInviteCode() (string, error) {
	code := make([]byte, codeLength)
	charsetLength := big.NewInt(int64(len(codeChars)))
	
	for i := 0; i < codeLength; i++ {
		// Generează un index aleatoriu securizat
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		
		// Selectează caracterul corespunzător
		code[i] = codeChars[randomIndex.Int64()]
	}
	
	return string(code), nil
}