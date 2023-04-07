package Parties

import (
	"bytes"
	"crypto/rand"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"math/big"
	"strconv"
)

type Verifier struct {
	parameters       Models.Parameters
	proverCommitment []byte
}

func (V Verifier) verifyOpening(triple Models.OpeningTriple) bool {
	position := triple.Index
	currentHash := triple.Value
	for _, value := range triple.OpenValues {
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
func (V Verifier) checkCorrectPebbleOfNode(tripleToCheck Models.OpeningTriple, parentTriples []Models.OpeningTriple) bool {
	//Get and check opening of the node itself
	if !V.verifyOpening(tripleToCheck) {
		return false
	}
	//Get and check all parents of the node
	parentHashes := make([]byte, 0, 1)
	for _, p := range parentTriples {
		if !V.verifyOpening(p) {
			return false
		}
		parentHashes = append(parentHashes, p.Value...)
	}

	//Compare to check that the node matches both the original graph and the merkle tree
	shouldBe := tripleToCheck.Value
	nodeLabel := []byte(strconv.Itoa(tripleToCheck.Index))
	toBeHashed := []byte(V.parameters.Id)
	toBeHashed = append(toBeHashed, nodeLabel...)
	toBeHashed = append(toBeHashed, parentHashes...)
	hash := hash_strategy.HashByteArray(toBeHashed)
	return bytes.Equal(hash, shouldBe)
}

func (V Verifier) InitializationPhase1(params Models.Parameters) {
	V.parameters = params
}

func (V Verifier) SaveCommitment(commitment []byte) {
	V.proverCommitment = commitment
}

func (V Verifier) PickChallenges() []int {
	size := V.parameters.GraphDescription.Size
	result := make([]int, size, size)
	for i := 0; i < size; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(size))) //Uniform random distribution of specific size, could maybe depend on parameters
		if err != nil {
			print("ERROR HAPPENED DURING SAMPLE")
			print(err.Error())
		}
		result[i] = int(randomIndex.Int64()) //If size of graph is less than 31 bits then this is fine.
	}
	return result
}

// parents := V.parameters.GraphDescription.GetParents(nodeIndex)
