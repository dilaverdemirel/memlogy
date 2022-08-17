package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDatabase() error {
	var err error

	db, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		return err
	}

	return db.Ping()
}

func Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS entry(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        log_time DATETIME NOT NULL UNIQUE,
        description TEXT NOT NULL
    );
    `

	_, err := db.Exec(query)
	return err
}
