package db

import (
	"app/models"
	"fmt"
)

func PollsCreateTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS single_proposal_polls (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			yes_votes INTEGER,
			no_votes INTEGER,
			resolved BOOLEAN
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS multiple_proposal_polls (
			id SERIAL PRIMARY KEY,
			proposal_ids INTEGER[]
		)
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateSingleProposalPoll(proposalID int) {
	_, err := db.Exec(`
		INSERT INTO single_proposal_polls (proposal_id, yes_votes, no_votes, resolved)
		VALUES ($1, 0, 0, FALSE)
	`, proposalID)
	if err != nil {
		fmt.Println(err)
	}
}

func GetUnResolvedPolls() []interface{} {
	rows, err := db.Query(`
		SELECT * FROM single_proposal_polls WHERE resolved = FALSE;
	`)
	if err != nil {
		fmt.Println(err)
	}

	polls := []interface{}{}
	for rows.Next() {
		p := struct {
			ID         int
			ProposalID int
			Proposal   models.Proposal
			YesVotes   int
			NoVotes    int
			Resolved   bool
		}{}
		err = rows.Scan(&p.ID, &p.ProposalID, &p.YesVotes, &p.NoVotes, &p.Resolved)
		if err != nil {
			fmt.Println(err)
		}
		p.Proposal = GetProposal(p.ProposalID)
		polls = append(polls, p)
	}
	rows.Close()

	return polls
}

func GetPoll(pollID int) models.SingleProposalPoll {
	poll := models.SingleProposalPoll{}

	fmt.Println("pollID", pollID)

	err := db.QueryRow(`
		SELECT * FROM single_proposal_polls WHERE id = $1;
	`, pollID).Scan(&poll.ID, &poll.ProposalID, &poll.YesVotes, &poll.NoVotes, &poll.Resolved)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("poll", poll)
	fmt.Println("poll.ProposalID", poll.ProposalID)

	poll.Proposal = GetProposal(poll.ProposalID)

	return poll
}

func VoteYes(pollID int) {
	_, err := db.Exec(`
		UPDATE single_proposal_polls
		SET yes_votes = yes_votes + 1
		WHERE id = $1
	`, pollID)
	if err != nil {
		fmt.Println(err)
	}
}

func VoteNo(pollID int) {
	_, err := db.Exec(`
		UPDATE single_proposal_polls
		SET no_votes = no_votes + 1
		WHERE id = $1
	`, pollID)
	if err != nil {
		fmt.Println(err)
	}
}
