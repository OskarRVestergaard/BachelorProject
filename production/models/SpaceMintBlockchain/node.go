package SpaceMintBlockchain

import (
	"bytes"
)

type node struct {
	block  Block
	length int
}

/*
returns 1 if the first node is greater

returns 0 if the nodes are equal

returns -1 if the second node is greater
*/

func (node1 *node) hasGreaterPathWeightThan(node2 node) int { //TODO CHANGE TO USE QUALITY (THERE SHOULD BE A QUALITY FUNCTION ON BLOCKS
	var lengthDifference = node1.length - node2.length
	if lengthDifference > 0 {
		return 1
	}
	if lengthDifference < 0 {
		return -1
	}

	//length is equal, therefore compare quality
	var node1quality, isGen1 = node1.block.GetQuality() // TODO THIS SHOULD NOT BE SINGLE BLOCK QUALITY, BUT BLOCK QUALITY
	var node2quality, isGen2 = node2.block.GetQuality() // TODO THIS SHOULD NOT BE SINGLE BLOCK QUALITY, BUT BLOCK QUALITY
	if isGen1 {
		return 1
	}
	if isGen2 {
		return -1
	}

	var byteComparison = bytes.Compare(node1quality.ToSlice(), node2quality.ToSlice()) //TODO This is actually fine as a quality function, given that every proof is the same size
	return byteComparison
}
