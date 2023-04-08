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

func (V *Verifier) verifyOpening(triple Models.OpeningTriple) bool {
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
func (V *Verifier) checkCorrectPebbleOfNode(tripleToCheck Models.OpeningTriple, parentTriples []Models.OpeningTriple) bool {
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

func (V *Verifier) InitializationPhase1(params Models.Parameters) {
	V.parameters = params
}

func (V *Verifier) SaveCommitment(commitment []byte) {
	V.proverCommitment = commitment
}

func (V *Verifier) PickChallenges() []int {
	size := V.parameters.GraphDescription.Size
	challengeAmount := size / 2
	result := make([]int, challengeAmount, challengeAmount)
	for i := 0; i < challengeAmount; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(size))) //Uniform random distribution of specific size, could maybe depend on parameters
		if err != nil {
			print("ERROR HAPPENED DURING SAMPLE")
			print(err.Error())
		}
		result[i] = int(randomIndex.Int64()) //If size of graph is less than 31 bits then this is fine.
	}
	return result
}

func (V *Verifier) VerifyChallenges(challenges []int, triples []Models.OpeningTriple, withGraphConsistencyCheck bool) bool {
	//TODO Add null checks (since this would indicate not being provided all needed information)

	//Organize triples
	size := len(challenges)
	tripleMap := make(map[int]Models.OpeningTriple, size)
	for _, triple := range triples {
		tripleMap[triple.Index] = triple
	}
	//Verify for each challenge that enough data was provided and that the data is correct.
	for _, challenge := range challenges {
		parents := V.parameters.GraphDescription.GetParents(challenge)
		challengeTriple := tripleMap[challenge]
		if withGraphConsistencyCheck {
			parentTriples := make([]Models.OpeningTriple, len(parents))
			for i, parentIndex := range parents {
				parentTriples[i] = tripleMap[parentIndex]
			}
			if !V.checkCorrectPebbleOfNode(challengeTriple, parentTriples) {
				return false
			}
		} else {
			if !V.verifyOpening(challengeTriple) {
				return false
			}
		}

	}
	return true
}
