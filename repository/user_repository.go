package repository

import (
	"Application/database"
	"Application/models"
	"database/sql"
	"errors"
	"strings"
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
			SELECT id, username, email, password
			FROM users
			WHERE username = ?
		`

	row := database.DB.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	query := `
			SELECT id, username, email, password
			FROM users
			WHERE email = ?
		`

	row := database.DB.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil

}

func CheckEmailExists(email string) (bool, error) {
	_, err := GetUserByEmail(email)
	if err != nil {
		if err.Error() == "User not found" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func CheckUsernameExists(username string) (bool, error) {
	_, err := GetUserByUsername(username)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if errMsg == "user not found" || errMsg == "user not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
