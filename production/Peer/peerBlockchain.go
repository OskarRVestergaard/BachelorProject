package Peer

import (
	"bytes"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"time"
)

func (p *Peer) createBlock(verificationKey string, slot int, draw lottery_strategy.WinningLotteryParams) blockchain.Block {
	//TODO Need to check that the draw is correct
	secretKey, foundSk := p.PublicToSecret[verificationKey]
	if !foundSk {
		panic("Tried to create a block but peer did not have the associated SecretKey")
	}
	p.blockTreeMutex.Lock()
	headBlock := p.blockTree.GetHead()
	headBlockHash := headBlock.HashOfBlock()
	transactionsToAdd := p.blockTree.GetTransactionsNotInTree(p.unfinalizedTransactions) //TODO If optimization is made, also do it here
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
		panic("Something went wrong, created block but gave it a wrong signature")
	}
}

func (p *Peer) SendBlockWithTransactions(slot int, draw lottery_strategy.WinningLotteryParams) {
	verificationKey := utils.GetSomeKey(p.PublicToSecret) //todo maybe make sure that it is the same public key that was used for the draw
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

func (p *Peer) verifyBlock(block blockchain.Block) bool {
	//TODO Needs to verify that the transactions are not already present too (just like the sender did), since someone not following the protocol could exploit this
	//TODO This is potentially very slow, but could be faster using dynamic programming in the case the chain best chain does not switch often
	if !block.HasCorrectSignature(p.signatureStrategy) {
		return false
	}
	if !p.verifyTransactions(block.BlockData.Transactions) {
		return false
	}
	if block.Draw.Vk == "DEBUG" { //TODO REMOVE THIS!
		return true
	}
	if block.Draw.Vk != block.Vk {
		return false
	}
	if !bytes.Equal(block.Draw.ParentHash, block.ParentHash) {
		return false
	}
	return p.lotteryStrategy.Verify(block.Vk, block.ParentHash, p.hardness, block.Draw.Counter)
}

func (p *Peer) verifyTransactions(transactions []blockchain.SignedTransaction) bool {
	for _, transaction := range transactions {
		transactionSignatureIsCorrect := utils.TransactionHasCorrectSignature(p.signatureStrategy, transaction)
		if !transactionSignatureIsCorrect {
			return false
		}
	}
	return true
}

func (p *Peer) handleBlock(block blockchain.Block) {
	if !p.verifyBlock(block) {
		return
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

func (p *Peer) StartMining() {
	verificationKey := utils.GetSomeKey(p.PublicToSecret)
	newHeadHashes := p.blockTree.SubScribeToGetHead()
	head := p.blockTree.GetHead()
	initialHash := head.HashOfBlock()
	winningDraws := make(chan lottery_strategy.WinningLotteryParams)
	p.lotteryStrategy.StartNewMiner(verificationKey, p.hardness, initialHash, newHeadHashes, winningDraws)
	go p.blockCreater(winningDraws)
}

func (p *Peer) blockCreater(wins chan lottery_strategy.WinningLotteryParams) {
	for {
		newWin := <-wins
		slot := p.blockTree.GetHead().Slot + 1
		p.SendBlockWithTransactions(slot, newWin)
	}
}
