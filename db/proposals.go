package db

import (
	"app/logic"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

func ProposalsCreateTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS proposals (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			submitted BOOLEAN DEFAULT FALSE,
			submitted_at TIMESTAMP,
			accepted BOOLEAN DEFAULT FALSE,
			accepted_at TIMESTAMP,
			rejected BOOLEAN DEFAULT FALSE,
			rejected_at TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %q", err)
	}
}

func CreateProposal(p logic.Proposal) int64 {
	_, err := db.Exec(`
		INSERT INTO proposals (title, description)
		VALUES ($1, $2);
	`, p.Title, p.Description)
	if err != nil {
		log.Fatalf("Error inserting proposal: %q", err)
	}

	return p.ID
}

func UpdateProposal(p logic.Proposal) {
	_, err := db.Exec(`
		UPDATE proposals
		SET title = $1, description = $2
		WHERE id = $3;
	`, p.Title, p.Description, p.ID)
	if err != nil {
		log.Fatalf("Error updating proposal: %q", err)
	}
}

func DeleteProposal(id int64) {
	_, err := db.Exec(`
		DELETE FROM proposals WHERE id = $1;
	`, id)
	if err != nil {
		log.Fatalf("Error deleting proposal: %q", err)
	}
}

func SubmitProposal(proposalId int64) {
	CreateSingleProposalPoll(proposalId)
	_, err := db.Exec(`
		UPDATE proposals
		SET submitted = TRUE, submitted_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`, proposalId)
	if err != nil {
		log.Fatalf("Error submitting proposal: %q", err)
	}
}

func GetProposals() []logic.Proposal {
	// Might be able to delete this
	rows, err := db.Query("SELECT id, title, description, created_at FROM proposals")
	if errors.Is(err, sql.ErrNoRows) {
		return []logic.Proposal{}
	}
	if err != nil {
		log.Fatalf("Error querying proposals: %q", err)
	}
	defer rows.Close()

	proposals := []logic.Proposal{}
	for rows.Next() {
		var p logic.Proposal
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

func GetUnsubmittedProposals() []logic.Proposal {
	rows, err := db.Query("SELECT id, title, description, created_at FROM proposals WHERE submitted = FALSE")
	if errors.Is(err, sql.ErrNoRows) {
		return []logic.Proposal{}
	}
	if err != nil {
		log.Fatalf("Error querying proposals: %q", err)
	}
	defer rows.Close()

	proposals := []logic.Proposal{}
	for rows.Next() {
		var p logic.Proposal
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

func GetProposal(id int64) logic.Proposal {
	var p logic.Proposal
	err := db.QueryRow(`
		SELECT id, title, description, created_at, submitted
		FROM proposals
		WHERE id = $1;
	`, id).Scan(&p.ID, &p.Title, &p.Description, &p.CreatedAt, &p.Submitted)
	if errors.Is(err, sql.ErrNoRows) {
		return logic.Proposal{}
	}
	if err != nil {
		log.Fatalf("Error querying proposal: %q", err)
	}

	changes := GetProposalChanges(p.ID)
	p.Changes = changes

	return p
}
