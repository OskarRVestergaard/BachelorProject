package models

import (
	"bytes"
	"strings"
)

type Blocktree struct {
	treeMap map[string]node
	head    node
}

func (tree *Blocktree) GetHead() Block {
	return tree.head.block
}

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
	var node1val = node1.block.GetVal()
	var node2val = node2.block.GetVal()
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

func (tree *Blocktree) HashToBlock(hash []byte) Block {
	return tree.treeMap[string(hash)].block
}

/*
AddBlock

returns 1 if successful.

returns 0 if the parent is not in the tree.

returns -1 is block is already in the tree.
*/
func (tree *Blocktree) AddBlock(block Block) int {

	//Check that this block is not already in the tree
	var newBlockHash = string(block.ToByteArray())
	var _, isAlreadyInTree = tree.treeMap[newBlockHash]
	if isAlreadyInTree {
		return -1
	}

	//Find parent
	var parentHash = block.H
	var parentNode, parentIsInTree = tree.treeMap[string(parentHash)]
	if !parentIsInTree {
		return 0
	}

	//Create and add the new block
	var newNode = node{
		block:  block,
		length: parentNode.length + 1,
	}
	tree.treeMap[newBlockHash] = newNode

	//Check if the longest chain has changed
	var newNodeGreater = newNode.hasGreaterPathWeightThan(tree.head)
	if newNodeGreater == 1 {
		tree.head = newNode
	}
	return 1
}
