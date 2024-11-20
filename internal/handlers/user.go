package handlers

import (
	"log"
	"net/http"
	"time"

	"app/internal/db"
	"app/internal/services"
)

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, r, "register.html.tmpl", nil, nil)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		err := services.RegisterUser(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, r, "login.html.tmpl", nil, nil)

	case http.MethodPost:

		log.Printf("Login request received")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Authenticate user
		user, err := services.AuthenticateUser(email, password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create a session for the user (we'll assume a session package or library is used here)
		sessionIDToken, err := db.CreateSession(user.ID, 30*24*time.Hour)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		// Set the session ID in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id_token",
			Value:    sessionIDToken,
			Expires:  time.Now().Add(time.Hour),
			HttpOnly: true,
			// TODO: Uncomment this for production to only send the cookie over HTTPS
			// Secure:   true,
			Path: "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the session (assuming the session library is in use)
	cookie, err := r.Cookie("session_id_token")
	if err == nil {
		/*
			http.SetCookie(w, &http.Cookie{
				Name:   "session_id_token",
				Value:  "",
				MaxAge: -1,
			})
		*/
		db.DeleteSession(cookie.Value)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
