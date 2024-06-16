package models

type Unit struct {
	ID       int
	Name     string
	Content  string
	GroupID  int
	Type     string // 'Existing', 'ProposedCreation', 'ProposedDeletion', 'ProposedModification', 'ProposedRename'
	ChangeID int
}

type Group struct {
	ID      int
	Name    string
	GroupID int
}

type Dependency struct {
	ID          int
	UnitID      int
	DependsOnID int
	Type        string // 'Existing', 'ProposedCreation', 'ProposedDeletion'
	ChangeID    int
}

type Graph struct {
	Units        []Unit
	Dependencies []Dependency
}

// PositionedUnits are used for rendering the graph

type PositionedUnit struct {
	Unit               Unit
	HorizontalPosition float64
}

type PositionedGraph struct {
	PositionedUnits []PositionedUnit
	Dependencies    []Dependency
}
