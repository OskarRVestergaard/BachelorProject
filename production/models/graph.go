package models

type Graph struct {
	// Assumed to be topologically sorted DAG according to index
	Size  int
	Edges [][]bool
	Value [][]byte
}
