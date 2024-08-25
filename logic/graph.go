package logic

import (
	"math"
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

	// For rendering
	HorizontalPosition float64
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

func (g *Graph) CalculateDirectAntecessorsAndSuccessors() {
	for i := range g.Units {
		g.Units[i].DirectAntecessors = nil // Resetting to avoid duplicating
		g.Units[i].DirectSuccessors = nil  // Resetting to avoid duplicating
	}

	for _, dependency := range g.Dependencies {
		for i := range g.Units {
			if (!dependency.UnitIsProposed && g.Units[i].ID == dependency.UnitID) ||
				(dependency.UnitIsProposed && g.Units[i].Type == "ProposedCreation" && g.Units[i].ChangeID == dependency.UnitID) {
				for j := range g.Units {
					if (!dependency.DependsOnIsProposed && g.Units[j].ID == dependency.DependsOnID) ||
						(dependency.DependsOnIsProposed && g.Units[j].Type == "ProposedCreation" && g.Units[j].ChangeID == dependency.DependsOnID) {
						g.Units[i].DirectAntecessors = append(g.Units[i].DirectAntecessors, &g.Units[j])
						g.Units[j].DirectSuccessors = append(g.Units[j].DirectSuccessors, &g.Units[i])
					}
				}
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

func (g *Graph) ToGonumGraph() *simple.DirectedGraph {
	directedGraph := simple.NewDirectedGraph()

	for _, unit := range g.Units {
		directedGraph.AddNode(Node{Unit: &unit})
	}

	for _, dependency := range g.Dependencies {
		// Find the from and to nodes
		fromNode := Node{}
		toNode := Node{}
		for _, unit := range g.Units {
			if (!dependency.DependsOnIsProposed && unit.ID == dependency.DependsOnID) ||
				(dependency.DependsOnIsProposed && unit.ChangeID == dependency.DependsOnID) {
				fromNode = Node{Unit: &unit}
				if toNode.Unit != nil {
					break
				}
			}
			if (!dependency.UnitIsProposed && unit.ID == dependency.UnitID) ||
				(dependency.UnitIsProposed && unit.ChangeID == dependency.UnitID) {
				toNode = Node{Unit: &unit}
				if fromNode.Unit != nil {
					break
				}
			}
		}
		directedGraph.SetEdge(simple.Edge{
			F: fromNode,
			T: toNode,
		})
	}

	return directedGraph
}

func (g *Graph) CalculateAccumulatedRelevances() {

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

func (g *Graph) SortAndPosition() {
	g.Sort()
	g.Position()
}

func (g *Graph) Sort() {
	// Vertical position is determined by the order of the nodes.
	// It is calculated by the topological sort algorithm, prioritizing by relevance.

	// No se si usar SortStabilized es lo que busco,
	// Primero agrupa las unidades por grupos conectados y luego hace lo que busco.
	// Los distintos grupos no están ordenador de forma inambigua.

	// Topological sort
	gonumGraph := g.ToGonumGraph()
	sortedNodes, err := topo.SortStabilized(gonumGraph, func(nodes []graph.Node) {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].(Node).Unit.AccumulatedRelevance > nodes[j].(Node).Unit.AccumulatedRelevance
		})
	})
	if err != nil {
		panic(err)
	}

	// Reorder g.Units according to sortedNodes
	newUnits := make([]Unit, len(sortedNodes))
	for i, node := range sortedNodes {
		newUnits[i] = *node.(Node).Unit
	}
	g.Units = newUnits
}

func (g *Graph) Position() {
	g.CalculateDirectAntecessorsAndSuccessors()
	g.InitializeHorizontalPositions()
	for i := 0; i < 10; i++ {
		g.EnsureMinimumDistanceOnBifurcations()
		g.AverageHorizontalPositionsBySuccessorsAndAntecessors()
		g.NormalizeHorizontalPositions()
		g.StraightenSingleConnections()
	}
	g.PutLonelyUnitsToZero()
}

func (g *Graph) InitializeHorizontalPositions() {
	// Initialize horizontal positions
	// TODO: This is just a shortcut
	for i := range g.Units {
		g.Units[i].HorizontalPosition = float64(i) / float64(len(g.Units)-1)
	}
}

func (g *Graph) AverageHorizontalPositionsBySuccessorsAndAntecessors() {
	// Average horizontal positions by successors and antecessors
	for j := range g.Units {
		unit := &g.Units[j]
		if len(unit.DirectSuccessors) > 1 {
			var average float64
			for _, successor := range unit.DirectSuccessors {
				average += successor.HorizontalPosition
			}
			unit.HorizontalPosition = average / float64(len(unit.DirectSuccessors))
		}
		if len(unit.DirectAntecessors) > 1 {
			var average float64
			for _, antecessor := range unit.DirectAntecessors {
				average += antecessor.HorizontalPosition
			}
			unit.HorizontalPosition = average / float64(len(unit.DirectAntecessors))
		}
	}
}

func (g *Graph) NormalizeHorizontalPositions() {
	// Normalize horizontal positions
	var min, max float64
	min = 1
	max = 0
	for _, unit := range g.Units {
		if unit.HorizontalPosition < min {
			min = unit.HorizontalPosition
		}
		if unit.HorizontalPosition > max {
			max = unit.HorizontalPosition
		}
	}
	for j := range g.Units {
		newPosition := (g.Units[j].HorizontalPosition - min) / (max - min)
		if math.IsNaN(newPosition) {
			newPosition = 0
		}
		g.Units[j].HorizontalPosition = newPosition
	}
}

func (g *Graph) StraightenSingleConnections() {
	// Straighten the single connections
	for i := 0; i < len(g.Units); i++ {
		unit := &g.Units[i]
		if !(len(unit.DirectSuccessors) > 1) &&
			(len(unit.DirectAntecessors) == 1 && len(unit.DirectAntecessors[0].DirectSuccessors) == 1) {
			unit.HorizontalPosition = unit.DirectAntecessors[0].HorizontalPosition
		}
	}
	for i := len(g.Units) - 1; i >= 0; i-- {
		unit := &g.Units[i]
		if !(len(unit.DirectAntecessors) > 1) &&
			(len(unit.DirectSuccessors) == 1 && len(unit.DirectSuccessors[0].DirectAntecessors) == 1) {
			unit.HorizontalPosition = unit.DirectSuccessors[0].HorizontalPosition
		}
	}
}

func (g *Graph) PutLonelyUnitsToZero() {
	// Put the lonely units to 0
	for i := range g.Units {
		unit := &g.Units[i]
		if len(unit.DirectSuccessors) == 0 && len(unit.DirectAntecessors) == 0 {
			unit.HorizontalPosition = 0
		}
	}
}

func (g *Graph) EnsureMinimumDistanceOnBifurcations() {
	minimumDistance := 0.2
	for i := 0; i < len(g.Units); i++ {
		unit := &g.Units[i]
		if len(unit.DirectAntecessors) > 1 && len(unit.DirectAntecessors) < int(1/minimumDistance) {
			var previous, current *Unit
			for j := 1; j < len(unit.DirectAntecessors); j++ {
				previous = unit.DirectAntecessors[j-1]
				current = unit.DirectAntecessors[j]
				if current.HorizontalPosition-previous.HorizontalPosition < minimumDistance {
					current.HorizontalPosition = previous.HorizontalPosition + minimumDistance
				}
			}
		}
	}
	for i := len(g.Units) - 1; i >= 0; i-- {
		unit := &g.Units[i]
		if len(unit.DirectSuccessors) > 1 && len(unit.DirectSuccessors) < int(1/minimumDistance) {
			var previous, current *Unit
			for j := 1; j < len(unit.DirectSuccessors); j++ {
				previous = unit.DirectSuccessors[j-1]
				current = unit.DirectSuccessors[j]
				if current.HorizontalPosition-previous.HorizontalPosition < minimumDistance {
					current.HorizontalPosition = previous.HorizontalPosition + minimumDistance
				}
			}
		}
	}
}
