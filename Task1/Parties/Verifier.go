package Parties

import (
	"crypto/rand"
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"math/big"
	"sort"
	"strconv"
)

type Verifier struct {
	parameters       PoSpaceModels.Parameters
	proverCommitment sha256.HashValue
}

func (V *Verifier) verifyOpening(triple PoSpaceModels.OpeningTriple) bool {
	position := triple.Index
	currentHash := triple.Value
	for _, value := range triple.OpenValues {
		isOdd := position%2 == 1
		if isOdd {
			currentHash = sha256.HashByteArray(append(value.ToSlice(), currentHash.ToSlice()...))
		} else {
			currentHash = sha256.HashByteArray(append(currentHash.ToSlice(), value.ToSlice()...))
		}
		position = position / 2
	}
	return currentHash.Equals(V.proverCommitment)
}

// CheckCorrectPebbleOfNode should be split since it uses information from "both sides" of the network traffic
func (V *Verifier) checkCorrectPebbleOfNode(tripleToCheck PoSpaceModels.OpeningTriple, parentTriples []PoSpaceModels.OpeningTriple) bool {
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
		parentHashes = append(parentHashes, p.Value.ToSlice()...)
	}

	//Compare to check that the node matches both the original graph and the merkle tree
	shouldBe := tripleToCheck.Value
	nodeLabel := []byte(strconv.Itoa(tripleToCheck.Index))
	toBeHashed := []byte(V.parameters.Id.String())
	toBeHashed = append(toBeHashed, nodeLabel...)
	toBeHashed = append(toBeHashed, parentHashes...)
	hash := sha256.HashByteArray(toBeHashed)
	result := hash.Equals(shouldBe)

	return result
}

func (V *Verifier) InitializationPhase1(params PoSpaceModels.Parameters) {
	V.parameters = params
}

func (V *Verifier) SaveCommitment(commitment sha256.HashValue) {
	V.proverCommitment = commitment
}

func (V *Verifier) PickChallenges() []int {
	size := V.parameters.GraphDescription.Size
	challengeAmount := 4
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

func (V *Verifier) VerifyChallenges(challenges []int, triples []PoSpaceModels.OpeningTriple, withGraphConsistencyCheck bool) bool {
	//Organize triples
	size := len(challenges)
	tripleMap := make(map[int]PoSpaceModels.OpeningTriple, size)
	for _, triple := range triples {
		tripleMap[triple.Index] = triple
	}
	usedCounter := 0
	for _, challenge := range challenges {
		challengeTriple, exists := tripleMap[challenge]
		if !exists {
			return false
		}
		usedCounter++
		parents := V.parameters.GraphDescription.GetPredecessors(challenge)
		sort.Slice(parents, func(i, j int) bool {
			return parents[i] < parents[j]
		})
		if withGraphConsistencyCheck {
			parentTriples := make([]PoSpaceModels.OpeningTriple, len(parents))
			for i, parentIndex := range parents {
				parentTriples[i], exists = tripleMap[parentIndex]
				if !exists {
					return false
				}
				usedCounter++
			}
			if !V.checkCorrectPebbleOfNode(challengeTriple, parentTriples) {
				return false
			}
		} else {
			if !V.verifyOpening(challengeTriple) {
				return false
			}
		}
		if usedCounter < len(tripleMap) {
			print("Prover sent everything it needed, but also send additional openings (which should not be allowed)")
		}
	}
	return true
}
