package models

// Units

type UnitCreation struct {
	ID         int64
	ProposalID int64
	Name       string
}

type UnitDeletion struct {
	ID         int64
	ProposalID int64
	UnitID     int64
}

type UnitRename struct {
	ID         int64
	ProposalID int64
	UnitID     int64
	Name       string
}

// Dependencies

type DependencyCreation struct {
	ID                  int64
	ProposalID          int64
	UnitIsProposed      bool
	UnitID              int64
	DependsOnIsProposed bool
	DependsOnID         int64
}

type DependencyDeletion struct {
	ID           int64
	ProposalID   int64
	DependencyID int64
}

// Content

type ContentModification struct {
	ID             int64
	ProposalID     int64
	UnitIsProposed bool
	UnitID         int64
	Content        string
}

type ContentFileUpload struct {
	ID             int64
	ProposalID     int64
	UnitIsProposed bool
	UnitID         int64
	// TODO
}

// Video

type VideoModification struct {
	ID         int64
	ProposalID int64
	UnitID     int64
	FromTime   int    // In miliseconds
	ToTime     int    // In miliseconds
	Content    string // Not a string
}

// Inherit Certifications/Read
// TODO
