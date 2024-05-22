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

func CreateProposal(p Proposal) {
	_, err := db.Exec(`
		INSERT INTO proposals (title, description)
		VALUES ($1, $2);
	`, p.Title, p.Description)
	if err != nil {
		log.Fatalf("Error inserting proposal: %q", err)
	}
}
