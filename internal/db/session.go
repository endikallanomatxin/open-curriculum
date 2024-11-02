package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"
)

func SessionsCreateTables() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		session_id_token TEXT NOT NULL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expiration_date TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID int, duration time.Duration) (string, error) {
	sessionIDToken, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	expirationDate := time.Now().Add(duration)
	_, err = db.Exec(`INSERT INTO sessions (session_id_token, user_id, expiration_date) VALUES ($1, $2, $3)`, sessionIDToken, userID, expirationDate)
	if err != nil {
		return "", err
	}

	return sessionIDToken, nil
}

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionExpired = errors.New("session expired")

func GetUserIDFromSession(sessionIDToken string) (int, error) {
	var userID int
	var expirationDate time.Time
	err := db.QueryRow(`SELECT user_id, expiration_date FROM sessions WHERE session_id_token = $1`, sessionIDToken).Scan(&userID, &expirationDate)
	if err == sql.ErrNoRows {
		return 0, ErrSessionNotFound
	}
	if err != nil {
		return 0, err
	}

	if expirationDate.Before(time.Now()) {
		DeleteSession(sessionIDToken)
		return 0, ErrSessionExpired
	}

	return userID, nil
}

func DeleteSession(sessionIDToken string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE session_id_token = $1`, sessionIDToken)
	return err
}
