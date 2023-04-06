package Task1

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

func GetDirectedAcyclicGraphStructure(id string) *Graph {
	size := 8
	edges := make([][]bool, size, size)
	for i := range edges {
		edges[i] = make([]bool, size, size)
	}
	edges[0][1] = true
	edges[0][2] = true
	edges[0][3] = true
	edges[1][3] = true
	edges[1][4] = true
	edges[2][4] = true
	edges[2][5] = true
	edges[0][7] = true
	edges[2][6] = true
	edges[3][6] = true
	edges[5][6] = true
	edges[5][7] = true

	resultGraph := &Graph{Size: size, Edges: edges, Value: make([][]byte, size, size)}

	return resultGraph
}
