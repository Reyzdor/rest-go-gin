package database

import "database/sql"

func createUsersTable(db *sql.DB) error {
	query :=
		`
 		CREATE TABLE IF NOT EXISTS users (
 		id INTEGER PRIMARY KEY AUTOINCREMENT,
 		username TEXT NOT NULL,
 		password TEXT NOT NULL,
 		email TEXT NOT NULL,
 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
 	);`

	_, err := db.Exec(query)
	return err
}
