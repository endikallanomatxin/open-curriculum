package models

type Unit struct {
	ID      int
	Name    string
	Content string
	GroupID int
}

type PositionedUnit struct {
	Unit               Unit
	HorizontalPosition float64
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
}

type Graph struct {
	Units        []Unit
	Dependencies []Dependency
}
