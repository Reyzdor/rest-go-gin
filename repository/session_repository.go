package repository

import (
	"database/sql"
	"errors"
	"time"
)

// Save user in DB

func SaveSession(db *sql.DB, userID int, refreshToken string) error {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	query := `
		INSERT INTO sessions (user_id, refresh_token, expires_at)
		VALUES (?, ?, ?)
	`

	_, err := db.Exec(query, userID, refreshToken, expiresAt)
	return err
}

// Get session on refresh token

func GetSession(db *sql.DB, refreshToken string) (int, error) {
	var userID int
	var expiresAt time.Time

	query := `
		SELECT user_id, expires_at FROM sessions
		WHERE refresh_token = ?
		LIMIT 1
	`

	err := db.QueryRow(query, refreshToken).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, nil
	}

	// Check validity period

	if time.Now().After(expiresAt) {
		db.Exec("DELETE FROM sessions WHERE refresh_token = ?", refreshToken)
		return 0, errors.New("Refresh token expired")
	}

	return userID, nil
}

// Delete session (exit)

func DeleteSession(db *sql.DB, refreshToken string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE refresh_token = ?", refreshToken)
	return err
}

// Delete all sessions user

func DeleteAllUserSessions(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
