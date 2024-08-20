package logic

import (
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type Unit struct {
	ID         int64
	Name       string
	Content    string
	GroupID    int64
	Type       string // 'Existing', 'ProposedCreation', 'ProposedDeletion', 'ProposedModification', 'ProposedRename'
	IsProposed bool
	ChangeID   int64

	Relevance float32 `default:"0.5"`

	// Later calculated
	AccumulatedRelevance float32
	DirectAntecessors    []*Unit
	DirectSuccessors     []*Unit
}

type Group struct {
	ID      int64
	Name    string
	GroupID int64
}

type Dependency struct {
	ID       int64
	Type     string // 'Existing', 'ProposedCreation', 'ProposedDeletion'
	ChangeID int64

	UnitIsProposed      bool
	UnitID              int64
	DependsOnIsProposed bool
	DependsOnID         int64
}

type Graph struct {
	Units        []Unit
	Dependencies []Dependency
}

func (g Graph) CalculateDirectAntecessorsAndSuccessors() {
	for _, dependency := range g.Dependencies {
		for j, unit := range g.Units {
			// Find the from and to units' indexes
			if (!dependency.UnitIsProposed && unit.ID == dependency.UnitID) ||
				(dependency.UnitIsProposed &&
					unit.Type == "ProposedCreation" &&
					unit.ChangeID == dependency.UnitID) {
				unit.DirectAntecessors = append(unit.DirectAntecessors, &g.Units[j])
			}
			if (!dependency.DependsOnIsProposed && unit.ID == dependency.DependsOnID) ||
				(dependency.DependsOnIsProposed &&
					unit.Type == "ProposedCreation" &&
					unit.ChangeID == dependency.DependsOnID) {
				unit.DirectSuccessors = append(unit.DirectSuccessors, &g.Units[j])
			}
		}
	}
}

type Node struct {
	Unit *Unit
}

func (n Node) ID() int64 {
	if n.Unit.Type == "Existing" ||
		n.Unit.Type == "ProposedDeletion" ||
		n.Unit.Type == "ProposedModification" ||
		n.Unit.Type == "ProposedRename" {
		return n.Unit.ID
	} else if n.Unit.Type == "ProposedCreation" {
		// To avoid id collisions, we use negative numbers for proposed creations
		return -n.Unit.ChangeID
	} else {
		panic("Invalid unit type")
	}
}

func (g Graph) ToGonumGraph() *simple.DirectedGraph {
	directedGraph := simple.NewDirectedGraph()

	for _, unit := range g.Units {
		directedGraph.AddNode(Node{Unit: &unit})
	}

	for _, dependency := range g.Dependencies {
		// Find the from and to nodes
		fromNode := Node{}
		toNode := Node{}
		for _, unit := range g.Units {
			if dependency.DependsOnIsProposed && unit.Type == "ProposedCreation" && unit.ChangeID == dependency.DependsOnID {
				fromNode = Node{Unit: &unit}
			} else if !dependency.DependsOnIsProposed && unit.ID == dependency.DependsOnID {
				fromNode = Node{Unit: &unit}
			}
			if dependency.UnitIsProposed && unit.Type == "ProposedCreation" && unit.ChangeID == dependency.UnitID {
				toNode = Node{Unit: &unit}
			} else if !dependency.UnitIsProposed && unit.ID == dependency.UnitID {
				toNode = Node{Unit: &unit}
			}

		}
		directedGraph.SetEdge(simple.Edge{
			F: fromNode,
			T: toNode,
		})
	}

	return directedGraph
}

func (g Graph) CalculateAccumulatedRelevances() {

	g.CalculateDirectAntecessorsAndSuccessors()

	gonumGraph := g.ToGonumGraph()

	orderedUnits, err := topo.Sort(gonumGraph)
	if err != nil {
		panic(err)
	}

	for i := len(orderedUnits) - 1; i >= 0; i-- {
		unit := orderedUnits[i].(Node).Unit
		unit.AccumulatedRelevance = unit.Relevance
		for _, antecessor := range unit.DirectAntecessors {
			unit.AccumulatedRelevance += antecessor.AccumulatedRelevance
		}
	}
}

// For graph rendering

type PositionedUnit struct {
	Unit               Unit
	HorizontalPosition float64
}

type PositionedGraph struct {
	PositionedUnits []PositionedUnit
	Dependencies    []Dependency
}

func (g Graph) Positioned() PositionedGraph {

	/*
		Vertical position. Order.
			Kahn algorithm for topological sort but prioritizing by relevance.

		Horizontal position.
			Find leaves order them by relevance and assign them a horizontal
			position (linspace from 0 to 1).
			For each node above in the graph, assign it the average horizontal
			position of its children.
	*/

	// Topological sort
	gonumGraph := g.ToGonumGraph()
	g.CalculateAccumulatedRelevances()
	sortedNodes, err := topo.SortStabilized(gonumGraph, func(nodes []graph.Node) {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].(Node).Unit.AccumulatedRelevance > nodes[j].(Node).Unit.AccumulatedRelevance
		})
	})
	// No se si usar SortStabilized es lo que busco,
	// Primero agrupa las unidades por grupos conectados y luego hace lo que busco.
	// Los distintos grupos no están ordenador de forma inambigua.
	if err != nil {
		panic(err)
	}

	// Order
	orderedUnits := []PositionedUnit{}
	for _, node := range sortedNodes {
		unit := *node.(Node).Unit
		orderedUnits = append(orderedUnits, PositionedUnit{
			Unit: unit,
		})
	}

	positionedGraph := PositionedGraph{
		PositionedUnits: orderedUnits,
		Dependencies:    g.Dependencies,
	}

	return positionedGraph
}
