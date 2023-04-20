package blockchain

import (
	"sync"
	"time"
)

/*
Blocktree

# A struct representing a Blocktree without any signature or transaction verification

Use the NewBlockTree method for creating a block tree!
*/
type Blocktree struct {
	treeMap                map[string]node
	head                   node
	subscriberChannelMutex sync.Mutex
	subscriberChannelList  []chan []byte
	newHeadBlocks          chan Block
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
	var genesisHash = genesisBlock.HashOfBlock()
	var genesisStringHash = string(genesisHash)
	treeMap[genesisStringHash] = genesisNode
	newHeadBlocks := make(chan Block)
	tree := &Blocktree{
		treeMap:       treeMap,
		head:          genesisNode,
		newHeadBlocks: newHeadBlocks,
	}
	tree.startSubscriptionHandler()
	return tree
}

/*
GetHead

returns the Block at the head of the longest currently known chain
*/
func (tree *Blocktree) GetHead() Block {
	return tree.head.block
}

/*
GetTransactionsNotInTree

returns the difference between transactions on block and list given
*/
func (tree *Blocktree) GetTransactionsNotInTree(unhandledTransactions []SignedTransaction) []SignedTransaction {

	head := tree.GetHead()
	transactionsInChain := tree.getTransactionsInChain(head)
	difference := getTransactionsInList1ButNotList2(unhandledTransactions, transactionsInChain)

	return difference
}

func (tree *Blocktree) getTransactionsInChain(block Block) []SignedTransaction {
	transactionsAccumulator := make([]SignedTransaction, 0)
	for !block.IsGenesis {
		transactionsAccumulator = append(transactionsAccumulator, block.BlockData.Transactions...)
		nextHash := block.ParentHash
		block = tree.HashToBlock(nextHash)
	}
	return transactionsAccumulator
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
	var newBlockHash = string(block.HashOfBlock())
	var _, isAlreadyInTree = tree.treeMap[newBlockHash]
	if isAlreadyInTree {
		return -1
	}

	//Find parent
	var parentHash = string(block.ParentHash)
	var parentNode, parentIsInTree = tree.treeMap[parentHash]
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
		tree.newHeadBlocks <- newNode.block
	}
	return 1
}

func (tree *Blocktree) startSubscriptionHandler() {
	tree.subscriberChannelMutex.Lock()
	tree.subscriberChannelList = make([]chan []byte, 0)
	go tree.subscriptionSubroutine()
	tree.subscriberChannelMutex.Unlock()
}

func (tree *Blocktree) subscriptionSubroutine() {
	for {
		newBlock := <-tree.newHeadBlocks
		tree.subscriberChannelMutex.Lock()
		for _, channel := range tree.subscriberChannelList {
			channel <- newBlock.HashOfBlock()
		}
		tree.subscriberChannelMutex.Unlock()
		time.Sleep(50 * time.Millisecond)
	}
}

func (tree *Blocktree) SubScribeToGetHead() (headHashes chan []byte) {
	newChannel := make(chan []byte)
	tree.subscriberChannelMutex.Lock()
	tree.subscriberChannelList = append(tree.subscriberChannelList, newChannel)
	tree.subscriberChannelMutex.Unlock()
	return newChannel
}
