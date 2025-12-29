package database

import "database/sql"

func createPostsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			main_image TEXT,
			price INT NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`

	query2 := `
		CREATE TABLE IF NOT EXISTS post_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			image_path TEXT NOT NULL,
			is_main BOOLEAN DEFAULT 0,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	_, err = db.Exec(query2)
	if err != nil {
		return err
	}

	return nil
}
