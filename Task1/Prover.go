package Task1

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"strconv"
)

func NewTestDAG(id string) *Graph {
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

	//TODO Used Known Random ID
	pebbleGraph(id, resultGraph)

	return resultGraph
}

func pebbleGraph(id string, graph *Graph) {
	// Assumed to be topologically sorted DAG according to index
	size := graph.Size
	for i := 0; i < size; i++ {
		vertexLabel := []byte(strconv.Itoa(i))
		toBeHashed := []byte(id)
		toBeHashed = append(toBeHashed, vertexLabel...)
		for j := 0; j < size; j++ {
			jIsParent := graph.Edges[j][i]
			if jIsParent {
				parentHashValue := graph.Value[j]
				toBeHashed = append(toBeHashed, parentHashValue...)
			}
		}

		//Debuggin' stuff
		//fmt.Println(i)
		//fmt.Println(vertexLabel)
		//fmt.Println(toBeHashed)

		graph.Value[i] = hash_strategy.HashByteArray(toBeHashed)
	}
}

func CreateMerkleTreeFromGraph(graph Graph) *MerkleTree {
	//Makes assumptions on the given graph, such as it being a DAG and sorted topologically by index
	size := graph.Size
	i := 1
	for i < size {
		i = i * 2
	}
	if i != size {
		panic("Graph must have 2^n number of nodes")
	}
	tree := MerkleTree{make([][]byte, size*2-1, size*2-1)}
	firstLeaf := size - 1
	//Inserting value for leaves
	for i := 0; i < size; i++ {
		tree.Nodes[firstLeaf+i] = graph.Value[i]
	}
	//Computing parents
	for i := firstLeaf - 1; i >= 0; i-- {
		leftChild := tree.Nodes[(i+1)*2-1]
		rightChild := tree.Nodes[(i+1)*2]
		toBeHashed := append(leftChild, rightChild...)
		tree.Nodes[i] = hash_strategy.HashByteArray(toBeHashed)
	}
	return &tree
}
