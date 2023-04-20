package blockchain

import (
	"bytes"
	"strings"
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
func (node1 *node) hasGreaterPathWeightThan(node2 node) int {
	var lengthDifference = node1.length - node2.length
	if lengthDifference > 0 {
		return 1
	}
	if lengthDifference < 0 {
		return -1
	}
	//length is equal compare val
	var node1val, isGenesis1 = node1.block.GetVal()
	if isGenesis1 {
		return 1
	}
	var node2val, isGenesis2 = node2.block.GetVal()
	if isGenesis2 {
		return -1
	}
	var stringComparison = strings.Compare(node1val, node2val)
	if stringComparison != 0 {
		return stringComparison
	}
	//Both length and val are equal
	//(some party send multiple blocks but with different data, or adding blocks to different chains of same length)
	//Therefore we sort the byte array of the block (should be fine since this is deterministic)
	var node1bytes = node1.block.ToByteArray()
	var node2bytes = node2.block.ToByteArray()
	var byteComparison = bytes.Compare(node1bytes, node2bytes)
	return byteComparison
}

//TODO Keep a list of children, such that cleanup is much... cleaner :)
