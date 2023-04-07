package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/pospace/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenVerification(t *testing.T) {
	id := "id"
	index := 4
	dag := graph.PrepareGraph(id)
	merkle := graph.CreateMerkleTreeFromGraph(*dag)
	oneOpen := merkle.Open(index)
	commitment := merkle.GetRootCommitment()
	openValue := merkle.GetLeaf(index)
	verification := graph.VerifyOpening(commitment, index, openValue, oneOpen)
	assert.True(t, verification)
}

func TestOpenVerificationWrongIndex(t *testing.T) {
	id := "id"
	index := 4
	dag := graph.PrepareGraph(id)
	merkle := graph.CreateMerkleTreeFromGraph(*dag)
	oneOpen := merkle.Open(index)
	commitment := merkle.GetRootCommitment()
	openValue := merkle.GetLeaf(index)
	verification := graph.VerifyOpening(commitment, 3, openValue, oneOpen)
	assert.False(t, verification)
}

func TestGetParents(t *testing.T) {
	id := "id"
	dag := graph.PrepareGraph(id)
	assert.Equal(t, 0, len(dag.GetParents(0)))
	assert.Equal(t, []int{0, 1}, dag.GetParents(3))
	assert.Equal(t, []int{2, 3, 5}, dag.GetParents(6))
}

func TestCheckPebblingOfNode(t *testing.T) {
	id := "id"
	index := 0
	dag := graph.PrepareGraph(id)
	merkle := graph.CreateMerkleTreeFromGraph(*dag)
	assert.True(t, graph.CheckCorrectPebbleOfNode(id, index, dag, merkle))
}
