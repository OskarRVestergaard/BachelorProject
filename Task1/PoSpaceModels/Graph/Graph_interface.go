package Graph

import "github.com/OskarRVestergaard/BachelorProject/production/sha256"

type Graph interface {
	InitGraph(n int, k int, withValues bool)
	AddEdge(from int, to int)
	IfEdge(from int, to int) bool
	GetSuccessors(node int) []int
	GetPredecessors(node int) []int
	GetSize() int
	GetValue() []sha256.HashValue
	SetValue([]sha256.HashValue)
	SortEdges()
	DebugInfo() (int, int, [][]int, [][]int, []sha256.HashValue)
}
