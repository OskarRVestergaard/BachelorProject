package Task1

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels"
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels/Graph"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/google/uuid"
	"math"
	mathRand "math/rand"
)

//Some slight modifications has been made to the original protocol, such as sending the index back and forth, the size
//of the graph and the distribution and amount of challenges.

//Maybe use int64 instead of switching between int types, and potentially allowing very big graphs

func generateDirectedAcyclicGraphStructure(size int) *Graph.Graph {
	resultGraph := &Graph.Graph{Size: size, Value: make([]sha256.HashValue, size, size)}
	resultGraph.InitGraph(size)
	resultGraph.AddEdge(0, 1)
	resultGraph.AddEdge(0, 2)
	resultGraph.AddEdge(0, 3)
	resultGraph.AddEdge(1, 3)
	resultGraph.AddEdge(1, 4)
	resultGraph.AddEdge(2, 4)
	resultGraph.AddEdge(2, 5)
	resultGraph.AddEdge(0, 7)
	resultGraph.AddEdge(2, 6)
	resultGraph.AddEdge(3, 6)
	resultGraph.AddEdge(5, 6)
	resultGraph.AddEdge(5, 7)

	return resultGraph
}

func GenerateTestingParameters() PoSpaceModels.Parameters {
	id := uuid.New()
	size := 8 //If changed, edge generation should also be made more general
	graphStructure := generateDirectedAcyclicGraphStructure(size)
	result := PoSpaceModels.Parameters{
		Id:               id,
		StorageBound:     2 * size,
		GraphDescription: graphStructure,
	}
	return result
}

func GenerateParameters(seed int64, n int, k int, alpha float64, beta float64, useForcedD bool, forcedD int) PoSpaceModels.Parameters {
	id := uuid.New()
	graphStructure := createGraph(seed, n, k, alpha, beta, useForcedD, forcedD)
	result := PoSpaceModels.Parameters{
		Id:               id,
		StorageBound:     2 * n * k,
		GraphDescription: graphStructure,
	}
	return result
}

func SimulateInitialization(prm PoSpaceModels.Parameters) (Parties.Prover, Parties.Verifier, bool) {
	prover := Parties.Prover{}
	verifier := Parties.Verifier{}

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

func createGraph(seed int64, n int, k int, alpha float64, beta float64, useForcedD bool, forcedD int) *Graph.Graph {
	size := n * k
	if !utils.PowerOfTwo(size) {
		panic("n and k must be a power of two")
	}
	graph := &Graph.Graph{Size: size, Value: make([]sha256.HashValue, size, size)}
	graph.InitGraph(size)
	var d int
	if useForcedD {
		d = forcedD
	} else {
		d = int(math.Ceil(CalculateD(alpha, beta)))
	}
	source := mathRand.NewSource(seed)
	rando := mathRand.New(source)

	preds := make([][]int, n, n)
	for i := range preds {
		preds[i] = make([]int, d, d)
		for k2 := range preds[i] {
			preds[i][k2] = -1
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
			graph.AddEdge(preds[i][j], n+i)
		}
	}

	for i := 0; i < k-2; i++ {
		firstIndex := (i + 1) * n
		for j := firstIndex; j < firstIndex+n; j++ {
			for y := firstIndex - n; y < firstIndex; y++ {
				if graph.IfEdge(y, j) {
					graph.AddEdge(y+n, j+n)
				}
			}
		}
	}

	return graph
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
