package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	// Obține calea către directorul de migrări
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("nu s-a putut obține directorul curent")
	}
	
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")
	
	// Creează driver-ul pentru postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("eroare la crearea driver-ului postgres: %v", err)
	}
	
	// Inițializează migrate
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("eroare la inițializarea migrate: %v", err)
	}
	
	// Execută migrările
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("eroare la executarea migrărilor: %v", err)
	}
	
	return nil
}