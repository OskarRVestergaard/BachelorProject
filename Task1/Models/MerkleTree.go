package Models

type MerkleTree struct {
	//Binary tree, children are at (index + 1) * 2 - 1 and (index + 1) * 2
	//Binary tree should be saturated fully
	Nodes [][]byte
}

func (tree *MerkleTree) GetRootCommitment() []byte {
	return tree.Nodes[0]
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
