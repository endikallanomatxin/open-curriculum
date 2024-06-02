package models

type SingleProposalPoll struct {
	ID         int
	ProposalID int
	Proposal   Proposal
	YesVotes   int
	NoVotes    int
	Resolved   bool `default:"false"`
}

type MultipleProposalPoll struct {
	ID          int
	ProposalIDs []int
}
