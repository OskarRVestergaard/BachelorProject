package SpaceMintBlockchain

import (
	"encoding/binary"
	"github.com/OskarRVestergaard/BachelorProject/Task1/PoSpaceModels"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy/PoSpace"
	"math/rand"
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
	genesisHash   sha256.HashValue
	nodeContainer map[sha256.HashValue]node
	head          node
	subscribers   chan []subscriber
	newHeadBlocks chan Block
	k             int
}

type subscriber struct {
	n               int
	miningLocations chan PoSpace.MiningLocation
}

/*
NewBlocktree

Constructor for making a new Blocktree, second parameter false if something went wrong such as the genesisBlock having IsGenesis equaling false
*/
func NewBlocktree(genesisBlock Block, k int) (Blocktree, bool) {
	if !genesisBlock.IsGenesis {
		return Blocktree{}, false
	}
	var blockQuality = models.CalculateQuality(genesisBlock.HashOfBlock(), 1000) //Does not really matter, but could be important if chain quality is fully implemented
	var genesisNode = node{
		block:              genesisBlock,
		length:             0,
		singleBlockQuality: blockQuality,
		chainQuality:       models.CalculateChainQuality([]float64{blockQuality}),
	}
	var genesisHash = genesisBlock.HashOfBlock()
	treeMapContainer := map[sha256.HashValue]node{}
	treeMapContainer[genesisHash] = genesisNode
	newHeadBlocks := make(chan Block, 20)
	tree := Blocktree{
		genesisHash:   genesisHash,
		nodeContainer: treeMapContainer,
		head:          genesisNode,
		newHeadBlocks: newHeadBlocks,
		k:             k,
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
	//Currently, since the lists are unsorted the algorithm just loops over all nm combinations, could be sorted first and then it would run in nlogn+mlogm
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
		nod, isEmpty := tree.hashToNode(nextHash, tree.nodeContainer)
		if isEmpty {
			panic("Parent in tree does not exist, tree is inconsistent")
		}
		block = nod.block
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
func (tree *Blocktree) HashToBlock(hash sha256.HashValue) (Block, bool) {
	nod, isEmpty := tree.hashToNode(hash, tree.nodeContainer)
	return nod.block, isEmpty
}

func (tree *Blocktree) hashToNode(hash sha256.HashValue, treeMap map[sha256.HashValue]node) (nod node, isEmpty bool) {
	result, foundKey := treeMap[hash]
	if !foundKey {
		return node{
			block:              Block{},
			length:             0,
			singleBlockQuality: 0,
			chainQuality:       0,
		}, true
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
func (tree *Blocktree) AddBlock(block Block, proofSizeN int64) int {

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
	var triplesHash = sha256.HashByteArray(PoSpaceModels.ListOfTripleToByteArray(block.HashSubBlock.Draw.ProofOfSpaceA))
	var blockQuality = models.CalculateQuality(triplesHash, proofSizeN)
	var chainQualities = append([]float64{blockQuality}, tree.collectBlockQualitiesForHead()...)
	var newNode = node{
		block:              block,
		length:             parentNode.length + 1,
		singleBlockQuality: blockQuality,
		chainQuality:       models.CalculateChainQuality(chainQualities),
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
	nod, isEmpty := tree.hashToNode(hashOfBlockToMineOn, tree.nodeContainer)
	if isEmpty {
		panic("GetMiningLocation called on invalid hash")
	}
	newBlock := nod.block
	challengeSetP, challengesSetV := tree.GetChallengesForExtendingOnBlockWithHash(hashOfBlockToMineOn, n*tree.k)
	newLocation := PoSpace.MiningLocation{
		Slot:          newBlock.TransactionSubBlock.Slot + 1, //This slot number is not used, the miner mine for every time slot
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
	if !reflect.DeepEqual(tree.head.block, comparisonTree.head.block) {
		return false //This conditional is not needed, but is here for debugging purposes
	}
	return true
}

func (tree *Blocktree) GetChallengesForExtendingOnHead(n int) (ProofChallengeSetP []int, CorrectCommitmentChallengesSetV []int) {
	//Should be calculated with dynamically fixed point prior in the chain according to the protocol as described on page 6
	head := tree.GetHead()
	return tree.GetChallengesForExtendingOnBlockWithHash(head.HashOfBlock(), n)
}

func (tree *Blocktree) GetChallengesForExtendingOnBlockWithHash(parentHash sha256.HashValue, n int) (ProofChallengeSetP []int, CorrectCommitmentChallengesSetV []int) {
	//TODO Should be calculated with dynamically fixed point prior in the chain according to the protocol as described on page 6 Here we just use the parent block
	challengesSamplingBlock, _ := tree.HashToBlock(parentHash)
	//We only sample from the proof chain to avoid some cases of challenge grinding
	hashSubBlockHash := sha256.HashByteArray(challengesSamplingBlock.HashSubBlock.ToByteArray())
	HashAsInt := int64(binary.LittleEndian.Uint64(hashSubBlockHash.ToSlice()))
	rnd := rand.New(rand.NewSource(HashAsInt)) // Math.rand is good for our case, since we want something deterministic given the seed and n

	challengeAmountA := 1 //TODO Discuss what this should be
	challengeAmountB := 2 //TODO Same as above

	challengesSetP := make([]int, challengeAmountA)
	challengesSetV := make([]int, challengeAmountB)
	for i := 0; i < challengeAmountA; i++ {
		challengeNumber := rnd.Intn(n)
		challengesSetP[i] = challengeNumber
	}
	for i := 0; i < challengeAmountB; i++ {
		challengeNumber := rnd.Intn(n)
		challengesSetV[i] = challengeNumber
	}
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
		nod, isEmpty := tree.hashToNode(nextHash, tree.nodeContainer)
		if isEmpty {
			panic("Parent in tree does not exist, tree is inconsistent")
		}
		currentNode = nod
		i++
		if i > 10000000 {
			panic("There is probably a cycle in what was supposed to be a tree")
		}
	}
	return qualityAccumulator
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
	currentNode, isEmpty := tree.hashToNode(blockHash, tree.nodeContainer)
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
	currentNode, isEmpty := tree.hashToNode(blockHash, tree.nodeContainer)
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
