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

	if err := createUsersTable(db); err != nil {
		return err
	}

	if err := createSessionsTable(db); err != nil {
		return err
	}

	if err := createPostsTable(db); err != nil {
		return err
	}

	DB = db
	return nil
}
