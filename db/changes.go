package db

// CHANGES
// This table links proposals with operations

type Change struct {
	ID            int
	ProposalID    int
	OperationType string // To avoid using needing to query each table
	OperationID   int
}

// OPERATIONS

// Units

type UnitCreation struct {
	ID     int
	UnitID int
	Name   string
}

type UnitDeletion struct {
	ID     int
	UnitID int
}

type UnitUpdate struct {
	ID     int
	UnitID int
	Name   string
}

// Dependencies

type DependencyCreation struct {
	ID          int
	UnitID      int
	DependsOnID int
}

type DependencyDeletion struct {
	ID          int
	UnitID      int
	DependsOnID int
}

// Documents
// (One per unit)

// Document Part Operations
// Analogous to a diff
// They are run from the end of the document to the beginning (so that line numbers don't change)

// Everything could be a change
// If it is a deletion, the content is an empty string
// If it is an addition, ToLine is same as FromLine (or maybe better: FromLine-1, pensarlo bien)

type DocumentModification struct {
	ID       int
	UnitID   int
	FromLine int
	ToLine   int
	Content  string // Can have multiple lines
}

type DocumentFileUpload struct {
	ID     int
	UnitID int
}

// Video
// (One per unit)

type VideoModification struct {
	ID       int
	UnitID   int
	FromTime int // In miliseconds
	ToTime   int // In miliseconds
	Content  string // Not a string
}

func ChangesCreateTables() {
}
