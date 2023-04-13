package blockchain

type Blocktree struct {
	treeMap map[string]node
	head    node
}

/*
GetHead

returns the Block at the head of the longest currently known chain
*/
func (tree *Blocktree) GetHead() Block {
	return tree.head.block
}

/*
HashToBlock

returns the Block that hashes to the parameter
*/
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
