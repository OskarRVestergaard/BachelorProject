package Peer

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"time"
)

// TODO FIX LATER
func (p *Peer) Mine() {
	var hasPotentialWinner bool
	for k := range p.PublicToSecret {
		hasPotentialWinner, _ = p.lotteryStrategy.Mine(k, "PrevHash")
	}

	if hasPotentialWinner {
		// do block stuff IDK
	}
}

func (p *Peer) createBlock(verificationKey string, slot int, draw string) blockchain.Block {
	//TODO Need to check that the draw is correct
	secretKey, foundSk := p.PublicToSecret[verificationKey]
	if !foundSk {
		panic("Tried to create a block but peer did not have the associated SecretKey")
	}
	p.blockTreeMutex.Lock()
	headBlock := p.blockTree.GetHead()
	headBlockHash := headBlock.HashOfBlock()
	transactionsToAdd := p.blockTree.GetTransactionsNotInTree(p.unfinalizedTransactions)
	resultBlock := blockchain.Block{
		IsGenesis: false,
		Vk:        verificationKey,
		Slot:      slot,
		Draw:      draw,
		BlockData: blockchain.BlockData{
			Transactions: transactionsToAdd,
		},
		ParentHash: headBlockHash,
		Signature:  nil,
	}
	resultBlock.SignBlock(p.signatureStrategy, secretKey)
	p.blockTreeMutex.Unlock()
	if resultBlock.HasCorrectSignature(p.signatureStrategy) {
		return resultBlock
	} else {
		panic("Created block but gave it a wrong signature")
	}
}

func (p *Peer) SendBlockWithTransactions(slot int, draw string) {
	verificationKey := utils.GetSomeKey(p.PublicToSecret)
	blockWithTransactions := p.createBlock(verificationKey, slot, draw)

	msg := blockchain.Message{
		MessageType:   constants.BlockDelivery,
		MessageSender: p.IpPort,
		MessageBlocks: []blockchain.Block{blockWithTransactions},
	}
	go p.FloodMessage(msg)
	for _, block := range msg.MessageBlocks {
		p.unhandledBlocks <- block
	}
}

func (p *Peer) startBlockHandler() {
	for {
		blockToHandle := <-p.unhandledBlocks
		go p.handleBlock(blockToHandle)
	}
}

func (p *Peer) handleBlock(block blockchain.Block) {
	//TODO The check are currently made here, this can hurt performance since some part might be done multiple times for a given block
	//TODO Needs to verify the draw
	//TODO Needs to verify that the transactions are not already present too (just like the sender did), since someone not following the protocol could exploit this
	blockSignatureIsCorrect := block.HasCorrectSignature(p.signatureStrategy)
	if !blockSignatureIsCorrect {
		return
	}
	//Check correctness of transactions
	transactions := block.BlockData.Transactions
	for _, transaction := range transactions {
		transactionSignatureIsCorrect := utils.TransactionHasCorrectSignature(p.signatureStrategy, transaction)
		if !transactionSignatureIsCorrect {
			return
		}
	}

	p.blockTreeMutex.Lock()
	var t = p.blockTree.AddBlock(block)
	switch t {
	case -2:
		//Block with isGenesis true, not a real block and should be ignored
	case -1:
		//Block is in tree already and can be ignored
	case 0:
		//Parent is not in the tree, try to add later
		//TODO Maybe have another slice that are blocks which are waiting for parents to be added,
		//TODO such that they can be added immediately follow the parents addition to the tree (in case 1)
		p.blockTreeMutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		p.blockTreeMutex.Lock()
		p.unhandledBlocks <- block
	case 1:
		//Block successfully added to the tree
	default:
		p.blockTreeMutex.Unlock()
		panic("addBlockReturnValueNotUnderstood")
	}
	p.blockTreeMutex.Unlock()
}

func (p *Peer) addTransaction(t blockchain.SignedTransaction) {
	p.unfinalizedTransMutex.Lock()
	p.unfinalizedTransactions = append(p.unfinalizedTransactions, t)
	p.unfinalizedTransMutex.Unlock()
}
