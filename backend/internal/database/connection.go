package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB,error) {
	db,err := sql.Open("postgres",databaseURL)
	if err != nil {
		return nil,fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w",err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrationFile := filepath.Join("internal","database","migrations","001_initial.sql")
	content,err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w",err)
	}

	if _,err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	return nil
}