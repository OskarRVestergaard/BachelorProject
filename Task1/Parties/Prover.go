package Parties

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"strconv"
)

type Prover struct {
	parameters   Models.Parameters
	pebbledGraph *Models.Graph
	merkleTree   *Models.MerkleTree
	commitment   sha256.HashValue
}

func (P *Prover) pebbleGraph() {
	// Assumed to be topologically sorted DAG according to index
	id := P.parameters.Id
	P.pebbledGraph = P.parameters.GraphDescription
	size := P.pebbledGraph.Size
	for i := 0; i < size; i++ {
		vertexLabel := []byte(strconv.Itoa(i))
		toBeHashed := []byte(id)
		toBeHashed = append(toBeHashed, vertexLabel...)
		for j := 0; j < size; j++ {
			jIsParent := P.pebbledGraph.Edges[j][i]
			if jIsParent {
				parentHashValue := P.pebbledGraph.Value[j].ToSlice()
				toBeHashed = append(toBeHashed, parentHashValue...)
			}
		}

		P.pebbledGraph.Value[i] = sha256.HashByteArray(toBeHashed)
	}
}

func (P *Prover) createMerkleTreeFromGraph() {
	//Makes assumptions on the given graph, such as it being a DAG and sorted topologically by index
	size := P.pebbledGraph.Size
	i := 1
	for i < size {
		i = i * 2
	}
	if i != size {
		panic("Graph must have 2^n number of nodes")
	}
	tree := Models.MerkleTree{Nodes: make([]sha256.HashValue, size*2-1, size*2-1)}
	firstLeaf := size - 1
	//Inserting value for leaves
	for i := 0; i < size; i++ {
		tree.Nodes[firstLeaf+i] = P.pebbledGraph.Value[i]
	}
	//Computing parents
	for i := firstLeaf - 1; i >= 0; i-- {
		leftChild := tree.Nodes[(i+1)*2-1].ToSlice()
		rightChild := tree.Nodes[(i+1)*2].ToSlice()
		toBeHashed := append(leftChild, rightChild...)
		tree.Nodes[i] = sha256.HashByteArray(toBeHashed)
	}
	P.merkleTree = &tree
	P.commitment = P.merkleTree.GetRootCommitment()
}

func (P *Prover) InitializationPhase1(params Models.Parameters) {
	P.parameters = params
	P.pebbleGraph()
	P.createMerkleTreeFromGraph()
}

func (P *Prover) GetCommitment() sha256.HashValue {
	return P.commitment
}

func (P *Prover) GetOpeningTriple(index int) (triple Models.OpeningTriple) {
	indexValue := P.merkleTree.GetLeaf(index)
	openingValues := P.merkleTree.Open(index)
	result := Models.OpeningTriple{
		Index:      index,
		Value:      indexValue,
		OpenValues: openingValues,
	}
	return result
}

func (P *Prover) AnswerChallenges(indices []int, withParents bool) (openingTriples []Models.OpeningTriple) {
	//Remove duplicates using a set
	var member struct{}
	indicesSet := make(map[int]struct{})
	for _, value := range indices {
		indicesSet[value] = member
	}
	//Find parents of the nodes
	if withParents {
		for index, _ := range indicesSet {
			parents := P.pebbledGraph.GetParents(index)
			for _, parent := range parents {
				indicesSet[parent] = member
			}
		}
	}
	//Append triple for each and return the result
	result := make([]Models.OpeningTriple, 0, 0)
	for i, _ := range indicesSet {
		result = append(result, P.GetOpeningTriple(i))
	}
	return result
}
