package logic

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

// Documents

type DocumentModification struct {
	// Analogous to a diff
	// They are run from the end of the document to the beginning (so that line numbers don't change)
	ID         int64
	ProposalID int64
	UnitID     int64
	FromLine   int
	ToLine     int
	Content    string // Can have multiple lines
	// If it is a deletion, the content is an empty string
	// If it is an addition, ToLine is same as FromLine (or maybe better: FromLine-1, pensarlo bien)
}

type DocumentFileUpload struct {
	ID         int64
	ProposalID int64
	UnitID     int64
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
