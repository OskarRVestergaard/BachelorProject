package PoWblockchain

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"reflect"
	"time"
)

/*
Blocktree

# A struct representing a Blocktree without any signature or transaction verification

Use the NewBlockTree method for creating a block tree!
The struct methods are NOT thread safe
*/
type Blocktree struct {
	genesisHash           sha256.HashValue
	nodeContainer         map[sha256.HashValue]node
	head                  node
	subscriberChannelList chan []chan sha256.HashValue
	newHeadBlocks         chan Block
}

/*
NewBlocktree

Constructor for making a new Blocktree, second parameter false if something went wrong such as the genesisBlock having IsGenesis equaling false
*/
func NewBlocktree(genesisBlock Block) (Blocktree, bool) {
	if !genesisBlock.IsGenesis {
		return Blocktree{}, false
	}
	var treeMap = map[sha256.HashValue]node{}
	var genesisNode = node{
		block:  genesisBlock,
		length: 0,
	}
	var genesisHash = genesisBlock.HashOfBlock()
	treeMap[genesisHash] = genesisNode
	newHeadBlocks := make(chan Block, 20)
	tree := Blocktree{
		genesisHash:   genesisHash,
		nodeContainer: treeMap,
		head:          genesisNode,
		newHeadBlocks: newHeadBlocks,
	}
	tree.startSubscriptionHandler()
	return tree, true
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
func (tree *Blocktree) GetTransactionsNotInTree(unhandledTransactions []models.SignedPaymentTransaction) []models.SignedPaymentTransaction {

	head := tree.GetHead()
	transactionsInChain := tree.getTransactionsInChain(head)
	difference := models.GetTransactionsInList1ButNotList2(unhandledTransactions, transactionsInChain)

	return difference
}

func (tree *Blocktree) getTransactionsInChain(block Block) []models.SignedPaymentTransaction {
	transactionsAccumulator := make([]models.SignedPaymentTransaction, 0)
	i := 0
	for !block.IsGenesis {
		transactionsAccumulator = append(transactionsAccumulator, block.BlockData.Transactions...)
		nextHash := block.ParentHash
		block = tree.HashToBlock(nextHash)
		i++
		if i > 10000000 {
			panic("There is probably a cycle in what was supposed to be a tree")
		}
	}
	return transactionsAccumulator
}

/*
HashToBlock

returns the Block that hashes to the parameter
*/
func (tree *Blocktree) HashToBlock(hash sha256.HashValue) Block {
	result, foundKey := tree.nodeContainer[hash]
	if !foundKey {
		panic("Hash given to tree is not in tree!")
	}
	return result.block
}

func (tree *Blocktree) hashToNode(hash sha256.HashValue) (nod node, isEmpty bool) {
	result, foundKey := tree.nodeContainer[hash]
	if !foundKey {
		return node{}, true
	}
	return result, false
}

/*
AddBlock

returns 1 if successful.

returns 0 if the parent is not in the tree.

returns -1 if block is already in the tree.

returns -2 if block is marked as genesis block

returns -3 if slot number is not greater than parent
*/
func (tree *Blocktree) AddBlock(block Block) int {

	//Refuse to add a new genesisBlock
	if block.IsGenesis {
		return -2
	}

	//Check that this block is not already in the tree
	var newBlockHash = block.HashOfBlock()
	var _, isAlreadyInTree = tree.nodeContainer[newBlockHash]
	if isAlreadyInTree {
		return -1
	}

	//Find parent
	var parentHash = block.ParentHash
	var parentNode, parentIsInTree = tree.nodeContainer[parentHash]
	if !parentIsInTree {
		return 0
	}

	//Check that slot is greater
	var newSlot = block.Slot
	var parentSlot = parentNode.block.Slot
	if newSlot <= parentSlot {
		return -3
	}

	//Create and add the new block
	var newNode = node{
		block:  block,
		length: parentNode.length + 1,
	}
	tree.nodeContainer[newBlockHash] = newNode

	//Check if the longest chain has changed
	var newNodeGreater = newNode.hasGreaterPathWeightThan(tree.head)
	if newNodeGreater == 1 {
		tree.head = newNode
		tree.newHeadBlocks <- newNode.block
	}
	return 1
}

func (tree *Blocktree) startSubscriptionHandler() {
	tree.subscriberChannelList = make(chan []chan sha256.HashValue, 1)
	subscriberChannelList := make([]chan sha256.HashValue, 0)
	tree.subscriberChannelList <- subscriberChannelList
	go tree.subscriptionSubroutine()
}

func (tree *Blocktree) subscriptionSubroutine() {
	for {
		newBlock := <-tree.newHeadBlocks
		subscriberChannelList := <-tree.subscriberChannelList
		for _, channel := range subscriberChannelList {
			go func(c chan sha256.HashValue) {
				c <- newBlock.HashOfBlock()
			}(channel)
		}
		tree.subscriberChannelList <- subscriberChannelList
		time.Sleep(50 * time.Millisecond)
	}
}

func (tree *Blocktree) SubScribeToGetHead() (headHashes chan sha256.HashValue) {
	newChannel := make(chan sha256.HashValue, 10)
	subscriberChannelList := <-tree.subscriberChannelList
	subscriberChannelList = append(subscriberChannelList, newChannel)
	tree.subscriberChannelList <- subscriberChannelList
	newChannel <- tree.head.block.HashOfBlock()
	return newChannel
}

func (tree *Blocktree) Equals(comparisonTree Blocktree) bool {
	//Is not thread safe, since the tree could change during operation
	if !reflect.DeepEqual(tree.nodeContainer, comparisonTree.nodeContainer) {
		return false
	}
	return reflect.DeepEqual(tree.head.block, comparisonTree.head.block)
}

//Tree visualization

/*
Method is probably very slow, since it checks the whole tree for children, going the other direction is much easier
It is only supposed to be used in testing
*/
func (tree *Blocktree) getChildren(parent node) []sha256.HashValue {
	var children []sha256.HashValue
	parentHash := parent.block.HashOfBlock()
	for hashValues, potentialChild := range tree.nodeContainer {
		if potentialChild.block.ParentHash.Equals(parentHash) {
			children = append(children, hashValues)
		}
	}
	return children
}

type visualNode struct {
	blockInformation node
	children         []visualNode
}

/*
HashToVisualNode
Probably very slow, to be used only for manual testing
Also not tail recursive
*/
func (tree *Blocktree) RootToVisualNode() visualNode {
	return tree.HashToVisualNode(tree.genesisHash)
}

func (tree *Blocktree) HashToVisualNode(blockHash sha256.HashValue) visualNode {
	currentNode, isEmpty := tree.hashToNode(blockHash)
	if isEmpty {
		panic("HashToVisualNode given hash not in use")
	}
	children := tree.getChildren(currentNode)
	var childrenVisualNodes []visualNode

	for _, child := range children {
		childrenVisualNodes = append(childrenVisualNodes, tree.HashToVisualNode(child))
	}

	return visualNode{
		blockInformation: currentNode,
		children:         childrenVisualNodes,
	}
}

type Chain struct {
	blockInformation node
	parent           *Chain
}

func (tree *Blocktree) HeadToChain() Chain {
	return tree.HashToChain(tree.head.block.HashOfBlock())
}

func (tree *Blocktree) HashToChain(blockHash sha256.HashValue) Chain {
	currentNode, isEmpty := tree.hashToNode(blockHash)
	if isEmpty {
		panic("HashToChain given hash not in use")
	}
	if currentNode.block.IsGenesis {
		return Chain{
			blockInformation: currentNode,
			parent:           nil,
		}
	}
	parentChain := tree.HashToChain(currentNode.block.ParentHash)
	return Chain{
		blockInformation: currentNode,
		parent:           &parentChain,
	}
}
