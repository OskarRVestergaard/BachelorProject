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
The struct methods are NOT thread safe (except for the handling of subscribers)
*/
type Blocktree struct {
	nodeContainer map[sha256.HashValue]node
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
	var blockQuality = CalculateQuality(genesisBlock)
	var genesisNode = node{
		block:              genesisBlock,
		length:             0,
		singleBlockQuality: blockQuality,
		chainQuality:       CalculateChainQuality([]float64{blockQuality}),
	}
	var genesisHash = genesisBlock.HashOfBlock()
	treeMapContainer := map[sha256.HashValue]node{}
	treeMapContainer[genesisHash] = genesisNode
	newHeadBlocks := make(chan Block, 20)
	tree := Blocktree{
		nodeContainer: treeMapContainer,
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
func (tree *Blocktree) GetTransactionsNotInTree(unhandledTransactions SpacemintTransactions) SpacemintTransactions {

	head := tree.GetHead()
	transactionsInChain := tree.getTransactionsInChain(head)
	transactionDifferences := getTransactionsNotInBlockChain(unhandledTransactions, transactionsInChain)

	return transactionDifferences
}

func getSpaceCommitsInList1ButNotList2(list1 []SpaceCommitment, list2 []SpaceCommitment) []SpaceCommitment {
	//Currently, since the lists are unsorted the algorithm just loops over all nm combinations, could be sorted first and then i would run in nlogn+mlogm
	var difference []SpaceCommitment
	found := false
	for _, val1 := range list1 {
		found = false
		for _, val2 := range list2 {
			if val1.Id == val2.Id {
				found = true
			}
		}
		if !found {
			difference = append(difference, val1)
		}
	}
	return difference
}

func getTransactionsNotInBlockChain(blockTransactions SpacemintTransactions, blockchainTransactions []SpacemintTransactions) SpacemintTransactions {
	var paymentAccumulator []models.SignedPaymentTransaction
	var spaceAccumulator []SpaceCommitment
	//var penaltyAccumulator []Penalty
	for _, currentBlockTransactions := range blockchainTransactions {
		paymentAccumulator = append(paymentAccumulator, currentBlockTransactions.Payments...)
		spaceAccumulator = append(spaceAccumulator, currentBlockTransactions.SpaceCommitments...)
		//penaltyAccumulator = append(penaltyAccumulator, currentBlockTransactions.Penalties...)
	}
	finalPayments := models.GetTransactionsInList1ButNotList2(blockTransactions.Payments, paymentAccumulator)
	finalSpaceCommits := getSpaceCommitsInList1ButNotList2(blockTransactions.SpaceCommitments, spaceAccumulator)
	//FinalPenalties...
	result := SpacemintTransactions{
		Payments:         finalPayments,
		SpaceCommitments: finalSpaceCommits,
		Penalties:        []Penalty{},
	}
	return result
}

func (tree *Blocktree) getTransactionsInChain(block Block) []SpacemintTransactions {
	transactionsAccumulator := make([]SpacemintTransactions, 0)
	i := 0
	for !block.IsGenesis {
		transactionsAccumulator = append(transactionsAccumulator, block.TransactionSubBlock.Transactions)
		nextHash := block.ParentHash
		block = tree.hashToNode(nextHash, tree.nodeContainer).block
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
	block := tree.hashToNode(hash, tree.nodeContainer).block
	return block
}

func (tree *Blocktree) hashToNode(hash sha256.HashValue, treeMap map[sha256.HashValue]node) node {
	result, foundKey := treeMap[hash]
	if !foundKey {
		panic("Hash given to tree is not in tree!")
	}
	return result
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
	var newSlot = block.TransactionSubBlock.Slot
	var parentSlot = parentNode.block.TransactionSubBlock.Slot
	if newSlot <= parentSlot {
		return -3
	}

	//Create and add the new block
	var blockQuality = CalculateQuality(block)
	var chainQualities = append([]float64{blockQuality}, tree.collectBlockQualitiesForHead()...)
	var newNode = node{
		block:              block,
		length:             parentNode.length + 1,
		singleBlockQuality: CalculateQuality(block),
		chainQuality:       CalculateChainQuality(chainQualities),
	}
	//Don't add node while subscribers are being notified
	subscribers := <-tree.subscribers
	tree.nodeContainer[newBlockHash] = newNode
	//Check if the longest chain has changed
	var newNodeGreater = newNode.hasGreaterPathWeightThan(tree.head) //TODO needs to not only use quality of single blocks, but compare the whole chain quality
	if newNodeGreater == 1 {
		tree.head = newNode
		tree.newHeadBlocks <- newNode.block
	}
	tree.subscribers <- subscribers
	return 1
}

func (tree *Blocktree) startSubscriptionHandler() {
	tree.subscribers = make(chan []subscriber, 1)
	subscribers := make([]subscriber, 0)
	tree.subscribers <- subscribers
	go tree.subscriptionSubroutine()
}

func (tree *Blocktree) GetMiningLocation(hashOfBlockToMineOn sha256.HashValue, n int) PoSpace.MiningLocation {
	newBlock := tree.hashToNode(hashOfBlockToMineOn, tree.nodeContainer).block
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
			func(sub subscriber) {
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
	treeMap1 := tree.nodeContainer
	treeMap2 := comparisonTree.nodeContainer
	if !reflect.DeepEqual(treeMap1, treeMap2) {
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
	//Also for efficiency, it would be better to split the finding random strings and getting the actual challenges
	//depending on n, into two parts, but this is only really important if multiple miners are active in the same tree,
	//which the code does not allow for anyways
	challengesSetP := []int{0, 1}
	challengesSetV := []int{0, 1, 2}
	return challengesSetP, challengesSetV
}

func (tree *Blocktree) collectBlockQualitiesForHead() (blockQualitiesFromHeadToGenesis []float64) {
	headNode := tree.head
	qualityAccumulator := make([]float64, headNode.length)
	currentNode := headNode

	i := 0
	for !currentNode.block.IsGenesis {
		qualityAccumulator = append(qualityAccumulator, currentNode.singleBlockQuality)
		nextHash := currentNode.block.ParentHash
		currentNode = tree.hashToNode(nextHash, tree.nodeContainer)
		i++
		if i > 10000000 {
			panic("There is probably a cycle in what was supposed to be a tree")
		}
	}
	return qualityAccumulator
}
