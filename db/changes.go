package db

type Change struct {
	ID          int
	OperationID int
}

// Units

type UnitCreation struct {
	ID     int
	UnitID int
}

type UnitDeletion struct {
	ID     int
	UnitID int
}

type UnitUpdate struct {
	ID          int
	UnitID      int
	Name        string
	Description string
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

type DocumentCreation struct {
	ID         int
	DocumentID int
}

type DocumentDeletion struct {
	ID         int
	DocumentID int
}

// Document Lines

type DocumentLineAddition struct {
	ID          int
	DocumentID  int
	Line        int
	LineContent string // Can have multiple lines
}

type DocumentLineDeletion struct {
	ID         int
	DocumentID int
	Line       int
}

type DocumentLineUpdate struct {
	ID          int
	DocumentID  int
	Line        int
	LineContent string
}

type DocumentLineMove struct {
	ID               int
	DocumentID       int
	Line             int
	TargetDocumentID int
	TargetLine       int
}

func ChangesCreateTables() {
}
