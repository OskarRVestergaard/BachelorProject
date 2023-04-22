package blockchain

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
	"log"
	"sync"
	"time"
)

/*
Blocktree

# A struct representing a Blocktree without any signature or transaction verification

Use the NewBlockTree method for creating a block tree!
*/
type Blocktree struct {
	treeMap                map[sha256.HashValue]node
	head                   node
	subscriberChannelMutex sync.Mutex
	subscriberChannelList  []chan []byte
	newHeadBlocks          chan Block
}

func byteSliceTo32ByteArray(bytes []byte) sha256.HashValue {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Tried to convert a byte slice to a hash-value but failed, probably because the slice had the wrong size!")
		}
	}()
	s4 := (*sha256.HashValue)(bytes)
	return *s4
}

/*
NewBlocktree

Constructor for making a new Blocktree, returns nil if isGenesis is false
*/
func NewBlocktree(genesisBlock Block) *Blocktree {
	if !genesisBlock.IsGenesis {
		return nil
	}
	var treeMap = map[sha256.HashValue]node{}
	var genesisNode = node{
		block:  genesisBlock,
		length: 0,
	}
	var genesisHash = genesisBlock.HashOfBlock()
	var genesisStringHash = byteSliceTo32ByteArray(genesisHash)
	treeMap[genesisStringHash] = genesisNode
	newHeadBlocks := make(chan Block, 20)
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
	i := 0
	for !block.IsGenesis {
		transactionsAccumulator = append(transactionsAccumulator, block.BlockData.Transactions...)
		nextHash := block.ParentHash
		block = tree.HashToBlock(nextHash)
		i++
		if i > 100 {
			panic("InfiniteLoop")
		}
	}
	return transactionsAccumulator
}

/*
HashToBlock

returns the Block that hashes to the parameter
*/
func (tree *Blocktree) HashToBlock(hash []byte) Block {
	result, foundKey := tree.treeMap[byteSliceTo32ByteArray(hash)]
	if !foundKey {
		panic("Hash given to tree is not in tree!")
	}
	return result.block
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
	var newBlockHash = byteSliceTo32ByteArray(block.HashOfBlock())
	var _, isAlreadyInTree = tree.treeMap[newBlockHash]
	if isAlreadyInTree {
		return -1
	}

	//Find parent
	var parentHash = byteSliceTo32ByteArray(block.ParentHash)
	var parentNode, parentIsInTree = tree.treeMap[parentHash]
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
			go func(c chan []byte) {
				c <- newBlock.HashOfBlock()
			}(channel)
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
