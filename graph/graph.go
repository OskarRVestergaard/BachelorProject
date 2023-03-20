package main

import (
	"example.com/packages/hash_strategy"
	"example.com/packages/models"
	"fmt"
	"strconv"
)

func main() {

	dag := newTestDAG()
	fmt.Println(*dag)
	merkle := CreateMerkleTree(*dag)
	fmt.Println(*merkle)
}

func newEmptyGraph(size int) *models.Graph {
	if size <= 0 {
		panic("Graph cannot have a size of 0 or less")
	} else {
		edges := make([][]bool, size, size)
		for i := range edges {
			edges[i] = make([]bool, size, size)
		}
		return &models.Graph{Size: size, Edges: edges, Value: make([][]byte, size, size)}
	}
}

func newTestDAG() *models.Graph {
	size := 6
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

	resultGraph := &models.Graph{Size: size, Edges: edges, Value: make([][]byte, size, size)}

	pebbleGraph(resultGraph)

	return resultGraph
}

func pebbleGraph(graph *models.Graph) {
	// Assumed to be topologically sorted DAG according to index
	size := graph.Size
	for i := 0; i < size; i++ {
		vertexLabel := []byte(strconv.Itoa(i))
		toBeHashed := vertexLabel
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

type MerkleTree struct {
	//Binary tree, children are at (index + 1) * 2 - 1 and (index + 1) * 2
	nodes [][]byte
}

func CreateMerkleTree(graph models.Graph) *MerkleTree {
	//TODO Assume we have a tree of size 2^N where N is a natural number
	//TODO Fake it, only uses the first 4 nodes of the graph
	size := 4 //TODO fake
	tree := MerkleTree{make([][]byte, size*2-1, size*2-1)}
	firstLeaf := size - 1
	//Inserting value for leaves
	for i := 0; i < size; i++ {
		tree.nodes[firstLeaf+i] = graph.Value[i]
	}
	for i := firstLeaf - 1; i >= 0; i-- {
		leftChild := tree.nodes[(i+1)*2-1]
		rightChild := tree.nodes[(i+1)*2]
		toBeHashed := append(leftChild, rightChild...)
		tree.nodes[i] = hash_strategy.HashByteArray(toBeHashed)
	}
	return &tree
}
