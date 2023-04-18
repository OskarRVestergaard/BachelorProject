package blockchain

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
	var genesisHash = genesisBlock.HashOfBlock()
	var genesisStringHash = string(genesisHash)
	treeMap[genesisStringHash] = genesisNode
	return &Blocktree{treeMap: treeMap, head: genesisNode}
}

/*
GetHead

returns the Block at the head of the longest currently known chain
*/
func (tree *Blocktree) GetHead() Block {
	return tree.head.block
}

/*
GetDifferenceInTransactions

returns the defference between transactions on block and list given
*/
func (tree *Blocktree) GetDifferenceInTransactions(unhandledTransactions []SignedTransaction) []SignedTransaction {
	var TransactionsOnBlockTree []SignedTransaction
	var TransOnBlock []SignedTransaction
	var TransactionsInBlock func(block Block) []SignedTransaction
	TransactionsInBlock = func(block Block) []SignedTransaction {
		if block.IsGenesis {
			return TransactionsOnBlockTree
		}
		TransOnBlock = append(TransOnBlock, block.BlockData.Transactions...)

		return TransactionsInBlock(tree.HashToBlock(block.ParentHash))
	}

	TransactionsInBlock(tree.GetHead())
	difference := differenceBetweenTwoTransactionLists(unhandledTransactions, TransOnBlock)
	print("asd")
	print(difference)

	return difference
}
func differenceBetweenTwoTransactionLists(list1 []SignedTransaction, list2 []SignedTransaction) []SignedTransaction {
	//https://www.tutorialspoint.com/golang-program-to-calculate-difference-between-two-slices
	var difference []SignedTransaction
	for _, val1 := range list1 {
		found := false
		for _, val2 := range list2 {
			if val1.Id == val2.Id {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, val1)
		}
	}
	return difference

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
	}
	return 1
}
