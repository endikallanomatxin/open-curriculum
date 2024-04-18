package db

type SingleProposalPoll struct {
	ID         int
	ProposalID int
}

type MultipleProposalPoll struct {
	ID int
}

func PollsCreateTables() {
	// Do nothing
}
