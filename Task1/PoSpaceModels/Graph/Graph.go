package Graph

import "github.com/OskarRVestergaard/BachelorProject/production/sha256"

type GeneralGraph struct {
	// Assumed to be topologically sorted DAG according to index
	size         int
	predecessors [][]int
	successors   [][]int
	value        []sha256.HashValue
}

func (graph *GeneralGraph) GetSize() int {
	return graph.size
}

func (graph *GeneralGraph) GetValue() []sha256.HashValue {
	return graph.value
}

func (graph *GeneralGraph) SetValue(value []sha256.HashValue) {
	graph.value = value
}

func (graph *GeneralGraph) InitGraph(n int, k int) {
	graph.size = n
	graph.value = make([]sha256.HashValue, n, n)

	graph.predecessors = make([][]int, n, n)
	for i := range graph.predecessors {
		graph.predecessors[i] = make([]int, 0)
	}

	graph.successors = make([][]int, n, n)
	for i := range graph.successors {
		graph.successors[i] = make([]int, 0)
	}

}

func (graph *GeneralGraph) AddEdge(from int, to int) {
	graph.successors[from] = append(graph.successors[from], to)
	graph.predecessors[to] = append(graph.predecessors[to], from)
}

func (graph *GeneralGraph) IfEdge(from int, to int) bool {
	result := false
	for _, i := range graph.predecessors[to] {
		if i == from {
			result = true
		}
	}
	return result
}

func (graph *GeneralGraph) GetSuccessors(node int) []int {
	return graph.successors[node]
}

// GetPredecessors returns the parents of a node
func (graph *GeneralGraph) GetPredecessors(node int) []int {
	return graph.predecessors[node]
}
