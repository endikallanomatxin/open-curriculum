package db

import (
	"log"

	_ "github.com/lib/pq"
)

// Proposal
type Proposal struct {
	// Collection of changes
	ID          int
	Title       string
	Description string
	CreatedAt   string
	Changes     []Change
}

func ProposalsCreateTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS proposals (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %q", err)
	}
}
