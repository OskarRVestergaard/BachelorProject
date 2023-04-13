package blockchain

import "github.com/OskarRVestergaard/BachelorProject/production/strategies/hash_strategy"

/*
Blocktree

# A struct representing a Blocktree without any signature or transaction verification

Use the NewBlockTree method for creating a block tree!
*/
type Blocktree struct {
	treeMap map[string]node
	head    node
}

/*
NewBlocktree

Constructor for making a new Blocktree, returns nil if isGenesis is false
*/
func NewBlocktree(genesisBlock Block) *Blocktree {
	if !genesisBlock.IsGenesis {
		return nil
	}
	var treeMap = map[string]node{}
	var genesisNode = node{
		block:  genesisBlock,
		length: 0,
	}
	var genesisHash = string(hash_strategy.HashByteArray(genesisBlock.ToByteArray()))
	treeMap[genesisHash] = genesisNode
	return &Blocktree{treeMap: treeMap}
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

returns -1 if block is already in the tree.

returns -2 if block is marked as genesis block
*/
func (tree *Blocktree) AddBlock(block Block) int {
	//Refuse to add a new genesisBlock
	if block.IsGenesis {
		return -2
	}

	//Check that this block is not already in the tree
	var newBlockHash = string(hash_strategy.HashByteArray(block.ToByteArray()))
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