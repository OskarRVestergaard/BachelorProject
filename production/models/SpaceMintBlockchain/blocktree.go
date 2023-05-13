package SpaceMintBlockchain

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy/PoSpace"
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
	treeMap       map[sha256.HashValue]node
	head          node
	subscribers   chan []subscriber
	newHeadBlocks chan Block
}

type subscriber struct {
	n               int
	miningLocations chan PoSpace.MiningLocation
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
		treeMap:       treeMap,
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
func (tree *Blocktree) GetTransactionsNotInTree(unhandledTransactions []models.SignedTransaction) []models.SignedTransaction {

	head := tree.GetHead()
	transactionsInChain := tree.getTransactionsInChain(head)
	difference := models.GetTransactionsInList1ButNotList2(unhandledTransactions, transactionsInChain)

	return difference
}

func (tree *Blocktree) getTransactionsInChain(block Block) []models.SignedTransaction {
	transactionsAccumulator := make([]models.SignedTransaction, 0)
	i := 0
	for !block.IsGenesis {
		transactionsAccumulator = append(transactionsAccumulator, block.TransactionSubBlock.Transactions.Payments...)
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
	result, foundKey := tree.treeMap[hash]
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
	var newBlockHash = block.HashOfBlock()
	var _, isAlreadyInTree = tree.treeMap[newBlockHash]
	if isAlreadyInTree {
		return -1
	}

	//Find parent
	var parentHash = block.ParentHash
	var parentNode, parentIsInTree = tree.treeMap[parentHash]
	if !parentIsInTree {
		return 0
	}

	//Check that slot is greater
	var newSlot = block.TransactionSubBlock.Slot
	var parentSlot = parentNode.block.TransactionSubBlock.Slot
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
	var newNodeGreater = newNode.hasGreaterPathWeightThan(tree.head) //TODO needs to not only use quality of single blocks, but compare the whole chain quality
	if newNodeGreater == 1 {
		tree.head = newNode
		tree.newHeadBlocks <- newNode.block
	}
	return 1
}

func (tree *Blocktree) startSubscriptionHandler() {
	tree.subscribers = make(chan []subscriber, 1)
	subscribers := make([]subscriber, 0)
	tree.subscribers <- subscribers
	go tree.subscriptionSubroutine()
}

func (tree *Blocktree) GetMiningLocation(hashOfBlockToMineOn sha256.HashValue, n int) PoSpace.MiningLocation {
	newBlock := tree.HashToBlock(hashOfBlockToMineOn)
	challengeSetP, challengesSetV := tree.getChallengesForExtendingOnBlockWithHash(hashOfBlockToMineOn, n)
	newLocation := PoSpace.MiningLocation{
		Slot:          newBlock.TransactionSubBlock.Slot + 1,
		ParentHash:    hashOfBlockToMineOn,
		ChallengeSetP: challengeSetP,
		ChallengeSetV: challengesSetV,
	}
	return newLocation
}

func (tree *Blocktree) subscriptionSubroutine() {
	for {
		newBlock := <-tree.newHeadBlocks
		hashOfBlockToExtendOn := newBlock.HashOfBlock()
		subscribers := <-tree.subscribers
		for _, singleSubscriber := range subscribers {
			go func(sub subscriber) {
				newLocation := tree.GetMiningLocation(hashOfBlockToExtendOn, sub.n)
				sub.miningLocations <- newLocation
			}(singleSubscriber)
		}
		tree.subscribers <- subscribers
		time.Sleep(50 * time.Millisecond)
	}
}

func (tree *Blocktree) SubScribeToGetHead(n int) (newHeadMiningLocations chan PoSpace.MiningLocation) {
	newMiningLocations := make(chan PoSpace.MiningLocation, 10)
	newSubscriber := subscriber{
		n:               n,
		miningLocations: newMiningLocations,
	}
	subscribers := <-tree.subscribers
	subscribers = append(subscribers, newSubscriber)
	tree.subscribers <- subscribers

	newMiningLocations <- tree.GetMiningLocation(tree.head.block.HashOfBlock(), n)
	return newMiningLocations
}

func (tree *Blocktree) Equals(comparisonTree Blocktree) bool {
	//Is not thread safe, since the tree could change during operation
	if !reflect.DeepEqual(tree.treeMap, comparisonTree.treeMap) {
		return false
	}
	return reflect.DeepEqual(tree.head.block, comparisonTree.head.block)
}

func (tree *Blocktree) GetChallengesForExtendingOnHead(n int) (ProofChallengeSetP []int, CorrectCommitmentChallengesSetV []int) {
	//Should be calculated with dynamically fixed point prior in the chain according to the protocol as described on page 6
	head := tree.GetHead()
	return tree.getChallengesForExtendingOnBlockWithHash(head.HashOfBlock(), n)
}

func (tree *Blocktree) getChallengesForExtendingOnBlockWithHash(parentHash sha256.HashValue, n int) (ProofChallengeSetP []int, CorrectCommitmentChallengesSetV []int) {
	//Should be calculated with dynamically fixed point prior in the chain according to the protocol as described on page 6
	//Todo change from fake challenges
	//Also for effeciency, it would be better to split the finding random strings and getting the actual challenges
	//depending on n, into two parts, but this is only really important if multiple miners are active in the same tree,
	//which the code does not allow for anyways
	challengesSetP := []int{0, 1}
	challengesSetV := []int{0, 1, 2}
	return challengesSetP, challengesSetV
}
