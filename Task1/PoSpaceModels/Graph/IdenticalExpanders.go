package Graph

import (
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"sort"
)

type IdenticalExpandersGraph struct {
	// Assumed to be topologically sorted DAG according to index
	n            int
	k            int
	predecessors [][]int
	successors   [][]int
	value        []sha256.HashValue
}

func (graph *IdenticalExpandersGraph) GetSize() int {
	return graph.n * graph.k
}

func (graph *IdenticalExpandersGraph) GetValue() []sha256.HashValue {
	return graph.value
}

func (graph *IdenticalExpandersGraph) SetValue(value []sha256.HashValue) {
	graph.value = value
}

func (graph *IdenticalExpandersGraph) InitGraph(n int, k int) {
	graph.n = n
	graph.k = k
	graph.value = make([]sha256.HashValue, n*k, n*k)

	graph.predecessors = make([][]int, n, n)
	for i := range graph.predecessors {
		graph.predecessors[i] = make([]int, 0)
	}

	graph.successors = make([][]int, n, n)
	for i := range graph.successors {
		graph.successors[i] = make([]int, 0)
	}
}

func (graph *IdenticalExpandersGraph) AddEdge(from int, to int) {
	if to/graph.n > 1 {
		panic("You can only add edges to the first bipartite expander!")
	}
	toto := to % graph.n
	graph.successors[from] = append(graph.successors[from], toto)
	graph.predecessors[toto] = append(graph.predecessors[toto], from)
}

func (graph *IdenticalExpandersGraph) IfEdge(from int, to int) bool {
	if from < graph.n {
		return false
	}
	result := false
	toto := to % graph.n
	fromfrom := from % graph.n
	for _, i := range graph.predecessors[toto] {
		if i == fromfrom {
			result = true
		}
	}
	return result
}

func (graph *IdenticalExpandersGraph) GetSuccessors(node int) []int {
	expanderNum := node / graph.n
	if expanderNum == graph.k-1 {
		return []int{}
	}
	nodeMod := node % graph.n
	successorsMod := graph.successors[nodeMod]
	successorsCopy := make([]int, len(successorsMod))
	copy(successorsCopy, successorsMod)
	for i, oldVal := range successorsCopy {
		successorsCopy[i] = oldVal + graph.n*(expanderNum+1)
	}
	return successorsCopy
}

func (graph *IdenticalExpandersGraph) GetPredecessors(node int) []int {
	expanderNum := node / graph.n
	if expanderNum == 0 {
		return []int{}
	}
	nodeMod := node % graph.n
	predecessorsMod := graph.predecessors[nodeMod]
	predecessorsCopy := make([]int, len(predecessorsMod))
	copy(predecessorsCopy, predecessorsMod)
	for i, oldVal := range predecessorsCopy {
		predecessorsCopy[i] = oldVal + graph.n*(expanderNum-1)
	}
	return predecessorsCopy
}

func (graph *IdenticalExpandersGraph) SortEdges() {
	for _, predecessors := range graph.predecessors {
		sort.Slice(predecessors, func(i, j int) bool {
			return predecessors[i] < predecessors[j]
		})
	}
	for _, successors := range graph.successors {
		sort.Slice(successors, func(i, j int) bool {
			return successors[i] < successors[j]
		})
	}
}
