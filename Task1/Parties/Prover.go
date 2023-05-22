package Parties

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels/Graph"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"sort"
	"strconv"

	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
)

type Prover struct {
	parameters   PoSpaceModels.Parameters
	pebbledGraph Graph.Graph
	merkleTree   *PoSpaceModels.MerkleTree
	commitment   sha256.HashValue
}

func (P *Prover) pebbleGraph() {
	// Assumed to be topologically sorted DAG according to index
	id := P.parameters.Id
	P.pebbledGraph = P.parameters.GraphDescription
	size := P.pebbledGraph.GetSize()
	tempHashingSlice := make([]byte, size*32+32) //The maximum size that can be hashed at any given point is if all nodes are involved, the constant factor should account for the label and id
	for i := 0; i < size; i++ {
		s := 0
		uniqueId := []byte(id.String())
		for _, b := range uniqueId {
			tempHashingSlice[s] = b
			s++
		}
		vertexLabel := []byte(strconv.Itoa(i))
		for _, b := range vertexLabel {
			tempHashingSlice[s] = b
			s++
		}
		parents := P.pebbledGraph.GetPredecessors(i)
		sort.Slice(parents, func(i, j int) bool {
			return parents[i] < parents[j]
		})
		for _, parent := range parents {
			parentHashValue := P.pebbledGraph.GetValue()[parent].ToSlice()
			for _, b := range parentHashValue {
				tempHashingSlice[s] = b
				s++
			}
		}
		values := P.pebbledGraph.GetValue()
		values[i] = sha256.HashByteArray(tempHashingSlice[0:s])
		P.pebbledGraph.SetValue(values)
	}
}

func (P *Prover) createMerkleTreeFromGraph() {
	//Makes assumptions on the given graph, such as it being a DAG and sorted topologically by index
	size := P.pebbledGraph.GetSize()

	if !utils.PowerOfTwo(size) {
		panic("Graph must have 2^n number of nodes")
	}
	tree := PoSpaceModels.MerkleTree{Nodes: make([]sha256.HashValue, size*2-1, size*2-1)}
	firstLeaf := size - 1
	//Inserting value for leaves
	for i := 0; i < size; i++ {
		tree.Nodes[firstLeaf+i] = P.pebbledGraph.GetValue()[i]
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

func (P *Prover) InitializationPhase1(params PoSpaceModels.Parameters) {
	P.parameters = params
	P.pebbleGraph()
	P.createMerkleTreeFromGraph()
}

func (P *Prover) GetCommitment() sha256.HashValue {
	return P.commitment
}

func (P *Prover) GetOpeningTriple(index int) (triple PoSpaceModels.OpeningTriple) {
	indexValue := P.merkleTree.GetLeaf(index)
	openingValues := P.merkleTree.Open(index)
	result := PoSpaceModels.OpeningTriple{
		Index:      index,
		Value:      indexValue,
		OpenValues: openingValues,
	}
	return result
}

func (P *Prover) AnswerChallenges(indices []int, withParents bool) (openingTriples []PoSpaceModels.OpeningTriple) {
	//Remove duplicates using a set
	var member struct{}
	challengeIndicesSet := make(map[int]struct{})
	for _, value := range indices {
		challengeIndicesSet[value] = member
	}
	//Find parents of the nodes
	parentIndicesSet := make(map[int]struct{})
	if withParents {
		for index, _ := range challengeIndicesSet {
			parents := P.pebbledGraph.GetPredecessors(index)
			for _, parent := range parents {
				parentIndicesSet[parent] = member
			}
		}
	}
	//Append triple for each and return the result
	result := make([]PoSpaceModels.OpeningTriple, 0, 0)
	for i, _ := range challengeIndicesSet {
		result = append(result, P.GetOpeningTriple(i))
	}
	for i, _ := range parentIndicesSet {
		result = append(result, P.GetOpeningTriple(i))
	}
	result = PoSpaceModels.SortOpeningTriples(result)
	return result
}
