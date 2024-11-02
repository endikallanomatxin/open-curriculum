package db

import (
	"database/sql"
	"log"

	"app/internal/models"
)

func UsersCreateTables() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		password_hash TEXT,
		oauth_provider TEXT NOT NULL,
		oauth_id TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(user *models.User) error {
	_, err := db.Exec(`INSERT INTO users (email, password_hash, oauth_provider, oauth_id)
		VALUES ($1, $2, $3, $4)`,
		user.Email, user.PasswordHash, user.OAuthProvider, user.OAuthID)
	return err
}

func GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(`SELECT id, email, password_hash, oauth_provider, oauth_id
		FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(`SELECT id, email, password_hash, oauth_provider, oauth_id
		FROM users WHERE email = $1`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func GetUserByOAuthID(provider, oauthID string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(`SELECT id, email, password_hash, oauth_provider, oauth_id
		FROM users WHERE oauth_provider = $1 AND oauth_id = $2`, provider, oauthID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func UpdateUser(user *models.User) error {
	_, err := db.Exec(`UPDATE users SET email = $1, password_hash = $2, oauth_provider = $3, oauth_id = $4
		WHERE id = $5`,
		user.Email, user.PasswordHash, user.OAuthProvider, user.OAuthID, user.ID)
	return err
}

func DeleteUser(id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
