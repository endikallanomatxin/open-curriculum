package db

import (
	"log"

	_ "github.com/lib/pq"

	models "app/models"
)

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

func CreateProposal(p models.Proposal) int {
	_, err := db.Exec(`
		INSERT INTO proposals (title, description)
		VALUES ($1, $2);
	`, p.Title, p.Description)
	if err != nil {
		log.Fatalf("Error inserting proposal: %q", err)
	}

	return p.ID
}

func UpdateProposal(p models.Proposal) {
	_, err := db.Exec(`
		UPDATE proposals
		SET title = $1, description = $2
		WHERE id = $3;
	`, p.Title, p.Description, p.ID)
	if err != nil {
		log.Fatalf("Error updating proposal: %q", err)
	}
}

func DeleteProposal(id int) {
	_, err := db.Exec(`
		DELETE FROM proposals WHERE id = $1;
	`, id)
	if err != nil {
		log.Fatalf("Error deleting proposal: %q", err)
	}
}

func GetProposals() []models.Proposal {
	rows, err := db.Query("SELECT id, title, description, created_at FROM proposals")
	if err != nil {
		log.Fatalf("Error querying proposals: %q", err)
	}
	defer rows.Close()

	proposals := []models.Proposal{}
	for rows.Next() {
		var p models.Proposal
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.CreatedAt)
		if err != nil {
			log.Fatalf("Error scanning proposal: %q", err)
		}
		changes := GetProposalChanges(p.ID)
		p.Changes = changes
		proposals = append(proposals, p)
	}

	return proposals
}

func GetProposal(id int) models.Proposal {
	var p models.Proposal
	err := db.QueryRow(`
		SELECT id, title, description, created_at
		FROM proposals
		WHERE id = $1;
	`, id).Scan(&p.ID, &p.Title, &p.Description, &p.CreatedAt)
	if err != nil {
		log.Fatalf("Error querying proposal: %q", err)
	}

	changes := GetProposalChanges(p.ID)
	p.Changes = changes

	return p
}
