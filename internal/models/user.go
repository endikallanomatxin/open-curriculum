package models

// NOTE:
// User implementation is done with the idea of supporting both
// - email/password users
// - OAuth users
// In the beginning, we only support email/password users,
// but we want to make it easy to add OAuth users later.

// TODO: Implement the OAuth user functionality

type User struct {
	ID            int
	Email         string
	PasswordHash  *string // Nullable, only set for email/password users
	OAuthProvider string  // e.g., "google" for OAuth users or "local" for email/password
	OAuthID       *string // Provider-specific ID for OAuth users
}
