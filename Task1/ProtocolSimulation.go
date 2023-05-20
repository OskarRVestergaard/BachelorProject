package Task1

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/google/uuid"
	"math"
	//"math/big"
	mathRand "math/rand"
)

//Some slight modifications has been made to the original protocol, such as sending the index back and forth, the size
//of the graph and the distribution and amount of challenges.

//Maybe use int64 instead of switching between int types, and potentially allowing very big graphs

func generateDirectedAcyclicGraphStructure(size int) *PoSpaceModels.Graph {
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
	edges[0][7] = true
	edges[2][6] = true
	edges[3][6] = true
	edges[5][6] = true
	edges[5][7] = true

	resultGraph := &PoSpaceModels.Graph{Size: size, Edges: edges, Value: make([]sha256.HashValue, size, size)}

	return resultGraph
}

func GenerateParameters(seed int64, n int, k int) PoSpaceModels.Parameters {
	id := uuid.New()
	graphEdges := createGraph(seed, n, k)
	result := PoSpaceModels.Parameters{
		Id:               id,
		StorageBound:     2 * n * k * n * k,
		GraphDescription: graphEdges,
	}
	return result
}

func SimulateInitialization() (Parties.Prover, Parties.Verifier, bool) {
	prover := Parties.Prover{}
	verifier := Parties.Verifier{}
	prm := GenerateParameters(5, 64, 64)

	//Prover and verifier gets ready for the protocol, the prover generates the hash-values of the graph
	//and computes a merkle tree commitment on it.
	prover.InitializationPhase1(prm)
	verifier.InitializationPhase1(prm)

	//Prover sends commitment to verifier
	rootMsg := prover.GetCommitment()
	verifier.SaveCommitment(rootMsg)

	//Verifier wants to proof the consistency of the given commitment and sends a set of challenges to the prover
	challenges := verifier.PickChallenges()

	//The prover sends all the openings of the challenges and their parents openings
	answerMsg := prover.AnswerChallenges(challenges, true)

	//Verifier checks that the openings are correct and that they match the hard to pebble graph
	result := verifier.VerifyChallenges(challenges, answerMsg, true)
	return prover, verifier, result
}

func SimulateExecution(prover Parties.Prover, verifier Parties.Verifier) bool {
	//The verifier picks challenges again
	challenges := verifier.PickChallenges()

	//The prover sends the openings of the challenges, but does not send the parent openings.
	answerMsg := prover.AnswerChallenges(challenges, false)

	//Verifier verifies the openings but does not consult the hard to pebble graph
	result := verifier.VerifyChallenges(challenges, answerMsg, false)
	return result
}

func createGraph(seed int64, n int, k int) *PoSpaceModels.Graph {
	if !utils.PowerOfTwo(n) {
		panic("n must be a power of two")
	}
	if !utils.PowerOfTwo(k) {
		panic("k must be a power of two")
	}

	edges := make([][]bool, n*k, n*k)
	for i := range edges {
		edges[i] = make([]bool, n*k, n*k)
	}

	var d = int(math.Ceil(CalculateD(0.25, 0.5)))
	source := mathRand.NewSource(seed)
	rando := mathRand.New(source)

	preds := make([][]int, n, n)
	for i := range preds {
		preds[i] = make([]int, d, d)
		for k := range preds[i] {
			preds[i][k] = -1
		}
		for j := 0; j < d; j++ {
			newNumber := false
			for !newNumber {
				random := rando.Intn(n)
				if !numberAlreadyChosen(random, preds[i]) {
					preds[i][j] = random
					newNumber = true
				}
			}
		}
	}

	for i := range preds {
		for j := range preds[i] {
			edges[preds[i][j]][n+i] = true
		}
	}

	for i := 0; i < len(edges)-n; i++ {
		for j := 0; j < len(edges)-n; j++ {
			if edges[i][j] { // == true
				edges[i+n][j+n] = true
			}
		}
	}

	resultGraph := &PoSpaceModels.Graph{Size: n * k, Edges: edges, Value: make([]sha256.HashValue, n*k, n*k)}

	return resultGraph
}

func numberAlreadyChosen(n int, lst []int) bool {
	for _, b := range lst {
		if b == n {
			return true
		}
	}
	return false
}

func CalculateD(alpha float64, beta float64) float64 {
	return (calculateEntropy(alpha) + calculateEntropy(beta)) / (calculateEntropy(alpha) - beta*calculateEntropy(alpha/beta))
}

func calculateEntropy(t float64) float64 {
	return -t*math.Log2(t)*t - (1-t)*math.Log2(1-t)
}
