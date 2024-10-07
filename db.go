package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize the SQLite database and create the table if it doesn't exist
func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %v", err)
	}

	// Create the table for URL mappings if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS url_mappings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_key TEXT NOT NULL UNIQUE,
		original_url TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

// Insert a new URL mapping into the database
func insertURLMapping(db *sql.DB, shortKey, originalURL string) error {
	insertQuery := `INSERT INTO url_mappings (short_key, original_url) VALUES (?, ?)`
	_, err := db.Exec(insertQuery, shortKey, originalURL)
	if err != nil {
		return fmt.Errorf("failed to insert URL mapping: %v", err)
	}
	return nil
}

// Retrieve the original URL by short key
func getOriginalURL(db *sql.DB, shortKey string) (string, error) {
	var originalURL string
	selectQuery := `SELECT original_url FROM url_mappings WHERE short_key = ? LIMIT 1`
	err := db.QueryRow(selectQuery, shortKey).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no URL found for short key: %s", shortKey)
	}
	if err != nil {
		return "", fmt.Errorf("failed to query URL: %v", err)
	}
	return originalURL, nil
}

// Find if a URL has already been shortened
func findURLMappingByOriginal(db *sql.DB, originalURL string) (string, error) {
	var shortKey string
	query := `SELECT short_key FROM url_mappings WHERE original_url = ?`
	err := db.QueryRow(query, originalURL).Scan(&shortKey)
	if err == sql.ErrNoRows {
		return "", nil // URL not found
	}
	return shortKey, err
}
