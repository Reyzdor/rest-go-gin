package repository

import (
	"Application/database"
	"Application/models"
	"database/sql"
	"errors"
)

func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)
	`

	_, err := database.DB.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password, email
		FROM users
		WHERE username = ?
	`

	row := database.DB.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
