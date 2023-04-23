package Peer

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"time"
)

func (p *Peer) createBlock(verificationKey string, slot int, draw lottery_strategy.WinningLotteryParams, blocktree blockchain.Blocktree) (newBlock blockchain.Block, isEmpty bool) {
	//TODO Need to check that the draw is correct
	secretKey, foundSk := p.PublicToSecret[verificationKey]
	if !foundSk {
		panic("Tried to create a block but peer did not have the associated SecretKey")
	}
	parentHash := draw.ParentHash
	p.unfinalizedTransMutex.Lock()
	allTransactionsToAdd := blocktree.GetTransactionsNotInTree(p.unfinalizedTransactions)
	var transactionsToAdd []blockchain.SignedTransaction
	if len(allTransactionsToAdd) == 0 {
		return blockchain.Block{}, true
	}
	if len(allTransactionsToAdd) < p.maximumTransactionsInBlock {
		transactionsToAdd = allTransactionsToAdd
	}
	if len(allTransactionsToAdd) > p.maximumTransactionsInBlock {
		transactionsToAdd = make([]blockchain.SignedTransaction, p.maximumTransactionsInBlock)
		for i := 0; i < p.maximumTransactionsInBlock; i++ {
			transactionsToAdd[i] = allTransactionsToAdd[i]
			//This could maybe cause starvation of transactions, if not enough blocks are made to saturate transaction demand
		}
	}
	p.unfinalizedTransMutex.Unlock()
	//
	resultBlock := blockchain.Block{
		IsGenesis: false,
		Vk:        verificationKey,
		Slot:      slot,
		Draw:      draw,
		BlockData: blockchain.BlockData{
			Transactions: transactionsToAdd,
		},
		ParentHash: parentHash,
		Signature:  nil,
	}
	resultBlock.SignBlock(p.signatureStrategy, secretKey)
	if resultBlock.HasCorrectSignature(p.signatureStrategy) {
		return resultBlock, false
	} else {
		panic("Something went wrong, created block but gave it a wrong signature")
	}
}

func (p *Peer) SendBlockWithTransactions(slot int, draw lottery_strategy.WinningLotteryParams) {
	verificationKey := utils.GetSomeKey(p.PublicToSecret) //todo maybe make sure that it is the same public key that was used for the draw
	blocktree := <-p.blockTreeChan
	blockWithTransactions, isEmpty := p.createBlock(verificationKey, slot, draw, blocktree)
	if isEmpty {
		p.blockTreeChan <- blocktree
		return
	}
	msg := blockchain.Message{
		MessageType:   constants.BlockDelivery,
		MessageSender: p.IpPort,
		MessageBlocks: []blockchain.Block{blockWithTransactions},
	}
	go p.FloodMessage(msg)
	for _, block := range msg.MessageBlocks {
		p.unhandledBlocks <- block
	}
	p.blockTreeChan <- blocktree
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
	if block.Draw.Vk != block.Vk {
		return false
	}
	if block.Draw.ParentHash != block.ParentHash {
		return false //TODO Instance of new block (slot2) being sent with an old draw (slot1)
	}
	if !p.lotteryStrategy.Verify(block.Vk, block.ParentHash, p.hardness, block.Draw.Counter) {
		return false
	}
	return true
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
	blocktree := <-p.blockTreeChan
	block = utils.MakeDeepCopyOfBlock(block)
	var t = blocktree.AddBlock(block)
	switch t {
	case -3:
		//Slot number is not greater than parent
		p.blockTreeChan <- blocktree
	case -2:
		//Block with isGenesis true, not a real block and should be ignored
		p.blockTreeChan <- blocktree
	case -1:
		//Block is in tree already and can be ignored
		p.blockTreeChan <- blocktree
	case 0:
		//Parent is not in the tree, try to add later
		//TODO Maybe have another slice that are blocks which are waiting for parents to be added,
		//TODO such that they can be added immediately follow the parents addition to the tree (in case 1)

		p.blockTreeChan <- blocktree
		time.Sleep(2000 * time.Millisecond) //Needs to be enough time for the other block to arrive
		p.unhandledBlocks <- block
	case 1:
		//Block successfully added to the tree
		p.blockTreeChan <- blocktree
	default:
		p.blockTreeChan <- blocktree
		panic("addBlockReturnValueNotUnderstood")
	}
}

func (p *Peer) addTransaction(t blockchain.SignedTransaction) {
	p.unfinalizedTransMutex.Lock()
	p.unfinalizedTransactions = append(p.unfinalizedTransactions, t)
	p.unfinalizedTransMutex.Unlock()
}

func (p *Peer) StartMining() {
	verificationKey := utils.GetSomeKey(p.PublicToSecret)
	blocktree := <-p.blockTreeChan
	newHeadHashes := blocktree.SubScribeToGetHead()
	head := blocktree.GetHead()
	initialHash := head.HashOfBlock()
	winningDraws := make(chan lottery_strategy.WinningLotteryParams)
	p.lotteryStrategy.StartNewMiner(verificationKey, p.hardness, initialHash, newHeadHashes, winningDraws)
	go p.blockCreater(winningDraws)

	p.blockTreeChan <- blocktree
}

func (p *Peer) blockCreater(wins chan lottery_strategy.WinningLotteryParams) {
	for {
		newWin := <-wins

		blocktree := <-p.blockTreeChan
		head := blocktree.GetHead()
		slot := head.Slot + 1
		p.blockTreeChan <- blocktree
		p.SendBlockWithTransactions(slot, newWin)
	}
}

func (p *Peer) GetBlockTree() blockchain.Blocktree {
	blocktree := <-p.blockTreeChan
	p.blockTreeChan <- blocktree
	return blocktree
}
