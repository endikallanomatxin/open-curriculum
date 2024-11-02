package models

type SingleProposalPoll struct {
	ID         int64
	ProposalID int64
	Proposal   Proposal
	YesVotes   int32
	NoVotes    int32
	Resolved   bool `default:"false"`
	Accepted   bool `default:"false"`
}

type MultipleProposalPoll struct {
	ID          int64
	ProposalIDs []int64
}
