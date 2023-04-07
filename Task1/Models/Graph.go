package Models

type Graph struct {
	// Assumed to be topologically sorted DAG according to index
	Size  int
	Edges [][]bool
	Value [][]byte
}

// GetParents returns the parents of a node, sorted low to high
func (graph *Graph) GetParents(nodeIndex int) []int {
	result := make([]int, 0, graph.Size)
	for j := 0; j < graph.Size; j++ {
		jIsParent := graph.Edges[j][nodeIndex]
		if jIsParent {
			result = append(result, j)
		}
	}
	return result
}
