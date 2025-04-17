package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config reprezintă configurația aplicației
type Config struct {
	// Server
	ServerPort  int
	Environment string
	AllowOrigins string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// Invite Code
	InviteCodeExpiration time.Duration
}

// LoadConfig încarcă configurația din variabilele de mediu
func LoadConfig() *Config {
	config := &Config{}

	// Server
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		port = 8080
	}
	config.ServerPort = port
	config.Environment = getEnv("ENVIRONMENT", "development")
	config.AllowOrigins = getEnv("ALLOW_ORIGINS", "*")

	// Database
	config.DatabaseURL = getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/relationship_helix?sslmode=disable")

	// JWT
	config.JWTSecret = getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")
	jwtExpiration, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		jwtExpiration = 24
	}
	config.JWTExpiration = time.Duration(jwtExpiration) * time.Hour

	// Invite Code
	inviteExpiration, err := strconv.Atoi(getEnv("INVITE_EXPIRATION_HOURS", "24"))
	if err != nil {
		inviteExpiration = 24
	}
	config.InviteCodeExpiration = time.Duration(inviteExpiration) * time.Hour

	return config
}

// getEnv obține o variabilă de mediu sau utilizează valoarea implicită dacă nu este setată
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return defaultValue
	}
	return value
}