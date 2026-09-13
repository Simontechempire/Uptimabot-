package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDatabase() error {
	var err error

	db, err = sql.Open("sqlite", "monitor.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS websites (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT UNIQUE NOT NULL,
			status TEXT NOT NULL DEFAULT 'CHECKING',
			response_time INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)

	return err
}

func saveWebsite(url string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO websites (url)
		VALUES (?)
	`, url)

	return err
}

func deleteWebsite(url string) error {
	_, err := db.Exec(`
		DELETE FROM websites
		WHERE url = ?
	`, url)

	return err
}

func updateDatabaseStatus(
	url string,
	status string,
	responseTime int64,
) error {
	_, err := db.Exec(`
		UPDATE websites
		SET status = ?, response_time = ?
		WHERE url = ?
	`, status, responseTime, url)

	return err
}
