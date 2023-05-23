package Graph

import "github.com/OskarRVestergaard/BachelorProject/production/sha256"

type Graph interface {
	InitGraph(n int, k int)
	AddEdge(from int, to int)
	IfEdge(from int, to int) bool
	GetSuccessors(node int) []int
	GetPredecessors(node int) []int
	GetSize() int
	GetValue() []sha256.HashValue
	SetValue([]sha256.HashValue)
	SortEdges()
}
