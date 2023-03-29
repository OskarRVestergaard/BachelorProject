package test

import (
	"github.com/OskarRVestergaard/BachelorProject/production/pospace/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenVerification(t *testing.T) {
	dag := graph.NewTestDAG()
	merkle := graph.CreateMerkleTree(*dag)
	oneOpen := merkle.Open(4)
	commitment := merkle.GetRootCommitment()
	openValue := merkle.Nodes[11]
	verification := graph.VerifyOpening(commitment, 4, openValue, oneOpen)
	assert.True(t, verification)
}

func TestOpenVerificationWrongIndex(t *testing.T) {
	dag := graph.NewTestDAG()
	merkle := graph.CreateMerkleTree(*dag)
	oneOpen := merkle.Open(4)
	commitment := merkle.GetRootCommitment()
	openValue := merkle.Nodes[11]
	verification := graph.VerifyOpening(commitment, 3, openValue, oneOpen)
	assert.False(t, verification)
}
