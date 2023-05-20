package SpaceMintBlockchain

type node struct {
	block              Block
	length             int
	singleBlockQuality float64
	chainQuality       float64
}

/*
returns 1 if the first node is greater

returns 0 if the nodes are equal

returns -1 if the second node is greater
*/

func (node1 *node) hasGreaterPathWeightThan(node2 node) int {
	var lengthDifference = node1.length - node2.length
	if lengthDifference > 0 {
		return 1
	}
	if lengthDifference < 0 {
		return -1
	}

	//length is equal, therefore compare quality
	var node1quality = node1.chainQuality
	var node2quality = node2.chainQuality

	if node1quality > node2quality {
		return 1
	}
	if node1quality < node2quality {
		return -1
	}
	return 0
}
