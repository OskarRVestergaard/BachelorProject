package PoSpaceModels

import "github.com/OskarRVestergaard/BachelorProject/production/sha256"

type Graph struct {
	// Assumed to be topologically sorted DAG according to index
	Size         int
	predecessors [][]int
	successors   [][]int
	Value        []sha256.HashValue
}

func (graph *Graph) InitGraph(size int) {
	graph.Size = size

	graph.predecessors = make([][]int, size, size)
	for i := range graph.predecessors {
		graph.predecessors[i] = make([]int, 0)
	}

	graph.successors = make([][]int, size, size)
	for i := range graph.successors {
		graph.successors[i] = make([]int, 0)
	}

}

func (graph *Graph) AddEdge(from int, to int) {
	graph.successors[from] = append(graph.successors[from], to)
	graph.predecessors[to] = append(graph.predecessors[to], from)
}

func (graph *Graph) IfEdge(from int, to int) bool {
	result := false
	for _, i := range graph.predecessors[to] {
		if i == from {
			result = true
		}
	}
	return result
}

func (graph *Graph) GetSuccessors(node int) []int {
	return graph.successors[node]
}

// GetPredecessors returns the parents of a node
func (graph *Graph) GetPredecessors(node int) []int {
	return graph.predecessors[node]
}
