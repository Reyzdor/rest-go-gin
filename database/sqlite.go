package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitSQLite() error {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(1)

	sc := `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(sc)
	if err != nil {
		return err
	}

	DB = db
	return nil
}
