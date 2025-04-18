package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// InitDB inițializează conexiunea la baza de date
func InitDB(databaseURL string) (*sql.DB, error) {
	// Conectare la baza de date
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("nu s-a putut deschide conexiunea cu baza de date: %v", err)
	}

	// Verifică conexiunea
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("nu s-a putut conecta la baza de date: %v", err)
	}

	// Setează parametrii de conexiune
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// RunMigrations execută migrările pentru baza de date
func RunMigrations(db *sql.DB) error {
	// Instead of using runtime.Caller, which relies on relative paths
	// Use an embedded approach or a fixed path that works in deployment

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("eroare la crearea driver-ului postgres: %v", err)
	}

	// Either use a fixed URL that points to migrations in your deployment
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations", // Update this path to where migrations exist in Docker
		"postgres", driver)

	if err != nil {
		return fmt.Errorf("eroare la inițializarea migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("eroare la executarea migrărilor: %v", err)
	}

	return nil
}
