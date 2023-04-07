package Task1

import (
	"crypto/rand"
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"math"
	"math/big"
)

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
