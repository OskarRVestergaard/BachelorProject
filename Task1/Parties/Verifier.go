package Parties

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"strconv"
)

type Verifier struct {
	parameters       Models.Parameters
	proverCommitment []byte
}

func (V Verifier) verifyOpening(openingIndex int, openingValue []byte, openingValues [][]byte) bool {
	position := openingIndex
	currentHash := openingValue
	for _, value := range openingValues {
		isOdd := position%2 == 1
		if isOdd {
			currentHash = hash_strategy.HashByteArray(append(value, currentHash...))
		} else {
			currentHash = hash_strategy.HashByteArray(append(currentHash, value...))
		}
		position = position / 2
	}
	return bytes.Equal(currentHash, V.proverCommitment)
}

// CheckCorrectPebbleOfNode should be split since it uses information from "both sides" of the network traffic
func (V Verifier) checkCorrectPebbleOfNode(id string, nodeIndex int, graph *Models.Graph, tree *Models.MerkleTree) bool {
	//Get and check opening of the node itself
	openingValue := tree.GetLeaf(nodeIndex)
	openingValues := tree.Open(nodeIndex)
	if !V.verifyOpening(nodeIndex, openingValue, openingValues) {
		return false
	}

	//Get and check all parents of the node
	parents := graph.GetParents(nodeIndex)
	parentHashes := make([]byte, 0, 1)
	for _, p := range parents {
		parentValue := tree.GetLeaf(p)
		parentOpeningValues := tree.Open(p)
		if !V.verifyOpening(p, parentValue, parentOpeningValues) {
			return false
		}
		parentHashes = append(parentHashes, parentValue...)
	}

	//Compare to check that the node matches both the original graph and the merkle tree
	shouldBe := openingValue
	nodeLabel := []byte(strconv.Itoa(nodeIndex))
	toBeHashed := []byte(id)
	toBeHashed = append(toBeHashed, nodeLabel...)
	toBeHashed = append(toBeHashed, parentHashes...)
	hash := hash_strategy.HashByteArray(toBeHashed)
	return bytes.Equal(hash, shouldBe)
}
