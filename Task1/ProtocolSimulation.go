package Task1

import (
	"crypto/rand"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Parties"
	"math"
	"math/big"
)

//Some slight modifications has been made to the original protocol, such as sending the index back and forth, the size
//of the graph and the distribution and amount of challenges.

//Maybe use int64 instead of switching between int types, and potentially allowing very big graphs

func generateDirectedAcyclicGraphStructure(size int) *Models.Graph {
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

	resultGraph := &Models.Graph{Size: size, Edges: edges, Value: make([][]byte, size, size)}

	return resultGraph
}

func generateParameters() Models.Parameters {
	random, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		print("ERROR HAPPENED:")
		print(err.Error())
	}
	id := random.String()
	size := 8 //If changed, edge generation should also be made more general
	graphEdges := generateDirectedAcyclicGraphStructure(size)
	result := Models.Parameters{
		Id:               id,
		StorageBound:     2 * size,
		GraphDescription: graphEdges,
	}
	return result
}

func SimulateInitialization() (Parties.Prover, Parties.Verifier, bool) {
	prover := Parties.Prover{}
	verifier := Parties.Verifier{}
	prm := generateParameters()

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
