package main

import (
	"example.com/packages/hash_strategy"
	"example.com/packages/models"
	"fmt"
	"strconv"
)

func main() {

	smallGraph := models.Graph{Size: 2, Edges: [][]bool{{false, true}, {false, false}}, Value: [][]byte{[]byte("one"), []byte("two")}}

	fmt.Println(smallGraph)
	fmt.Println(*newEmptyGraph(4))
	fmt.Println(*newTestDAG())

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
