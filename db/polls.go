package db

import (
	"app/logic"
	"fmt"
)

func PollsCreateTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS single_proposal_polls (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			yes_votes INTEGER DEFAULT 0,
			no_votes INTEGER DEFAULT 0,
			resolved BOOLEAN DEFAULT FALSE,
			accepted BOOLEAN DEFAULT FALSE
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

func CreateSingleProposalPoll(proposalID int64) {
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
		p := logic.SingleProposalPoll{}
		err = rows.Scan(&p.ID, &p.ProposalID, &p.YesVotes, &p.NoVotes, &p.Resolved, &p.Accepted)
		if err != nil {
			fmt.Println(err)
		}
		p.Proposal = GetProposal(p.ProposalID)
		polls = append(polls, p)
	}
	rows.Close()

	return polls
}

func GetAcceptedPolls() []interface{} {
	rows, err := db.Query(`
		SELECT * FROM single_proposal_polls WHERE accepted = TRUE;
	`)
	if err != nil {
		fmt.Println(err)
	}

	polls := []interface{}{}

	for rows.Next() {
		p := logic.SingleProposalPoll{}
		err = rows.Scan(&p.ID, &p.ProposalID, &p.YesVotes, &p.NoVotes, &p.Resolved, &p.Accepted)
		if err != nil {
			fmt.Println(err)
		}
		p.Proposal = GetProposal(p.ProposalID)
		polls = append(polls, p)
	}
	rows.Close()

	return polls
}

func GetPoll(pollID int64) logic.SingleProposalPoll {
	poll := logic.SingleProposalPoll{}

	err := db.QueryRow(`
		SELECT * FROM single_proposal_polls WHERE id = $1;
	`, pollID).Scan(&poll.ID, &poll.ProposalID, &poll.YesVotes, &poll.NoVotes, &poll.Resolved, &poll.Accepted)
	if err != nil {
		fmt.Println(err)
	}

	poll.Proposal = GetProposal(poll.ProposalID)

	return poll
}

func CheckPoll(pollID int64) {
	// If yes-no is greater than 10, resolve the poll
	poll := GetPoll(pollID)
	if poll.YesVotes-poll.NoVotes > 10 {
		_, err := db.Exec(`
			UPDATE single_proposal_polls
			SET resolved = TRUE, accepted = TRUE
			WHERE id = $1
		`, pollID)
		if err != nil {
			fmt.Println(err)
		}
		UpdateGraph()
		return
	}
	if poll.YesVotes-poll.NoVotes < -10 {
		_, err := db.Exec(`
			UPDATE single_proposal_polls
			SET resolved = TRUE, accepted = FALSE
			WHERE id = $1
		`, pollID)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}

func VoteYes(pollID int64) {
	_, err := db.Exec(`
		UPDATE single_proposal_polls
		SET yes_votes = yes_votes + 1
		WHERE id = $1
	`, pollID)
	if err != nil {
		fmt.Println(err)
	}
	CheckPoll(pollID)
}

func VoteNo(pollID int64) {
	_, err := db.Exec(`
		UPDATE single_proposal_polls
		SET no_votes = no_votes + 1
		WHERE id = $1
	`, pollID)
	if err != nil {
		fmt.Println(err)
	}
	CheckPoll(pollID)
}
