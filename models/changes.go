package models

// Units

type UnitCreation struct {
	ID         int
	ProposalID int
	Name       string
}

type UnitDeletion struct {
	ID         int
	ProposalID int
	UnitID     int
}

type UnitRename struct {
	ID         int
	ProposalID int
	UnitID     int
	Name       string
}

// Dependencies

type DependencyCreation struct {
	ID                  int
	ProposalID          int
	UnitIsProposed      bool
	UnitID              int
	DependsOnIsProposed bool
	DependsOnID         int
}

type DependencyDeletion struct {
	ID           int
	ProposalID   int
	DependencyID int
}

// Documents

type DocumentModification struct {
	// Analogous to a diff
	// They are run from the end of the document to the beginning (so that line numbers don't change)
	ID         int
	ProposalID int
	UnitID     int
	FromLine   int
	ToLine     int
	Content    string // Can have multiple lines
	// If it is a deletion, the content is an empty string
	// If it is an addition, ToLine is same as FromLine (or maybe better: FromLine-1, pensarlo bien)
}

type DocumentFileUpload struct {
	ID         int
	ProposalID int
	UnitID     int
	// TODO
}

// Video

type VideoModification struct {
	ID         int
	ProposalID int
	UnitID     int
	FromTime   int    // In miliseconds
	ToTime     int    // In miliseconds
	Content    string // Not a string
}

// Inherit Certifications/Read
// TODO
