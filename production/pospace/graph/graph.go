package graph

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"
	"strconv"
)

type Graph struct {
	// Assumed to be topologically sorted DAG according to index
	Size  int
	Edges [][]bool
	Value [][]byte
}

// GetParents returns the parents of a node, sorted low to high
func (graph *Graph) GetParents(nodeIndex int) []int {
	result := make([]int, 0, graph.Size)
	for j := 0; j < graph.Size; j++ {
		jIsParent := graph.Edges[j][nodeIndex]
		if jIsParent {
			result = append(result, j)
		}
	}
	return result
}

func newEmptyGraph(size int) *Graph {
	if size <= 0 {
		panic("Graph cannot have a size of 0 or less")
	} else {
		edges := make([][]bool, size, size)
		for i := range edges {
			edges[i] = make([]bool, size, size)
		}
		return &Graph{Size: size, Edges: edges, Value: make([][]byte, size, size)}
	}
}

func NewTestDAG(id string) *Graph {
	size := 8
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

	resultGraph := &Graph{Size: size, Edges: edges, Value: make([][]byte, size, size)}

	//TODO Used Known Random ID
	pebbleGraph(id, resultGraph)

	return resultGraph
}

func pebbleGraph(id string, graph *Graph) {
	// Assumed to be topologically sorted DAG according to index
	size := graph.Size
	for i := 0; i < size; i++ {
		vertexLabel := []byte(strconv.Itoa(i))
		toBeHashed := []byte(id)
		toBeHashed = append(toBeHashed, vertexLabel...)
		for j := 0; j < size; j++ {
			jIsParent := graph.Edges[j][i]
			if jIsParent {
				parentHashValue := graph.Value[j]
				toBeHashed = append(toBeHashed, parentHashValue...)
			}
		}

		//Debuggin' stuff
		//fmt.Println(i)
		//fmt.Println(vertexLabel)
		//fmt.Println(toBeHashed)

		graph.Value[i] = hash_strategy.HashByteArray(toBeHashed)
	}
}

type MerkleTree struct {
	//Binary tree, children are at (index + 1) * 2 - 1 and (index + 1) * 2
	Nodes [][]byte
}

func (tree *MerkleTree) GetRootCommitment() []byte {
	return tree.Nodes[0]
}

func CreateMerkleTree(graph Graph) *MerkleTree {
	size := graph.Size
	i := 1
	for i < size {
		i = i * 2
	}
	if i != size {
		panic("Graph must have 2^n number of nodes")
	}
	tree := MerkleTree{make([][]byte, size*2-1, size*2-1)}
	firstLeaf := size - 1
	//Inserting value for leaves
	for i := 0; i < size; i++ {
		tree.Nodes[firstLeaf+i] = graph.Value[i]
	}
	//Computing parents
	for i := firstLeaf - 1; i >= 0; i-- {
		leftChild := tree.Nodes[(i+1)*2-1]
		rightChild := tree.Nodes[(i+1)*2]
		toBeHashed := append(leftChild, rightChild...)
		tree.Nodes[i] = hash_strategy.HashByteArray(toBeHashed)
	}
	return &tree
}

func (tree *MerkleTree) GetLeaf(leafIndex int) []byte {
	firstLeaf := len(tree.Nodes) / 2
	indexInMerkleTree := firstLeaf + leafIndex
	result := tree.Nodes[indexInMerkleTree]
	return result
}

func (tree *MerkleTree) Open(openingIndex int) [][]byte {
	if openingIndex < 0 {
		panic("Index in merkle tree to open must not be negative!")
	}
	result := make([][]byte, 0, 1) //maybe instead of 1 choose math.Log(float64(len(tree.nodes))) (maximum size of nodes used in opening) THIS IS JUST AN OPTIMIZATION
	var isEven bool
	firstLeaf := len(tree.Nodes) / 2
	//some loop
	i := openingIndex + firstLeaf
	j := 0
	for i > 0 {
		isEven = (i-firstLeaf)%2 == 0
		if isEven {
			j = i + 1
		} else {
			j = i - 1
		}
		result = append(result, tree.Nodes[j])
		i = (i+1)/2 - 1 //Go to parent
	}
	return result
}

func VerifyOpening(commitment []byte, openingIndex int, openingValue []byte, openingValues [][]byte) bool {
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
	return bytes.Equal(currentHash, commitment)
}

// CheckCorrectPebbleOfNode should be split since it uses information from "both sides" of the network traffic
func CheckCorrectPebbleOfNode(id string, nodeIndex int, graph *Graph, tree *MerkleTree) bool {
	//Get and check opening of the node itself
	commitment := tree.GetRootCommitment()
	openingValue := tree.GetLeaf(nodeIndex)
	openingValues := tree.Open(nodeIndex)
	if !VerifyOpening(commitment, nodeIndex, openingValue, openingValues) {
		return false
	}

	//Get and check all parents of the node
	parents := graph.GetParents(nodeIndex)
	parentHashes := make([]byte, 0, 1)
	for _, p := range parents {
		parentValue := tree.GetLeaf(p)
		parentOpeningValues := tree.Open(p)
		if !VerifyOpening(commitment, p, parentValue, parentOpeningValues) {
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
