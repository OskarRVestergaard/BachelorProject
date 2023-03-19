package main

import "fmt"

func main() {

	smallGraph := Graph{edges: [][]bool{{false, true}, {false, false}}, value: [][]byte{[]byte("one"), []byte("two")}}

	fmt.Println(smallGraph)
	fmt.Println(*newEmptyGraph(4))
	fmt.Println(*newTestDAG())

}

type Graph struct {
	// Assumed to be topologically sorted DAG according to index
	edges [][]bool
	value [][]byte
}

func newEmptyGraph(size int) *Graph {
	if size <= 0 {
		panic("Graph cannot have a size of 0 or less")
	} else {
		edges := make([][]bool, size, size)
		for i := range edges {
			edges[i] = make([]bool, size, size)
		}
		return &Graph{edges, [][]byte{}}
	}
}

func newTestDAG() *Graph {
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

	return &Graph{edges, [][]byte{}}
}
