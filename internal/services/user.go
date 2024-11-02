package services

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"app/internal/db"
	"app/internal/models"
)

func HashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPasswordBytes), err
}

func RegisterUser(email, password string) error {
	existingUser, err := db.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:         email,
		PasswordHash:  &hashedPassword,
		OAuthProvider: "local",
	}

	err = db.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func AuthenticateUser(email, password string) (*models.User, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

var ErrSessionCookieNotFound = errors.New("session cookie not found")

func GetSessionUser(r *http.Request) (*models.User, error) {
	sessionIDTokenCookie, err := r.Cookie("session_id_token")
	if err != nil {
		return nil, ErrSessionCookieNotFound
	}

	sessionIDToken := sessionIDTokenCookie.Value

	userID, err := db.GetUserIDFromSession(sessionIDToken)
	if err != nil {
		return nil, err
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
