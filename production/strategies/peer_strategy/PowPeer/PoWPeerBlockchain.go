package PowPeer

import (
	"errors"
	"github.com/OskarRVestergaard/BachelorProject/production/Message"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"time"
)

func (p *PoWPeer) GetBlockTree() interface{} {
	//TODO Not at all thread safe to use it this way, fine if used for reading during testing
	blocktree := <-p.blockTreeChan
	p.blockTreeChan <- blocktree
	return blocktree
}

func (p *PoWPeer) StartMining() error {
	noActiveMiner := p.isMiningMutex.TryLock()
	if !noActiveMiner {
		return errors.New("peer is already mining")
	}
	secretKeys := <-p.publicToSecret
	verificationKey := utils.GetSomeKey(secretKeys)
	p.publicToSecret <- secretKeys
	blocktree := <-p.blockTreeChan
	newHeadHashes := blocktree.SubScribeToGetHead()
	head := blocktree.GetHead()
	initialHash := head.HashOfBlock()
	winningDraws := make(chan lottery_strategy.WinningLotteryParams, 10)
	p.lotteryStrategy.StartNewMiner(verificationKey, p.hardness, initialHash, newHeadHashes, winningDraws, p.stopMiningSignal)
	go p.blockCreatingLoop(winningDraws)

	p.blockTreeChan <- blocktree
	return nil
}

func (p *PoWPeer) StopMining() error {
	noActiveMiner := p.isMiningMutex.TryLock()
	if noActiveMiner {
		p.isMiningMutex.Unlock()
		return errors.New("peer is already not mining")
	}
	p.stopMiningSignal <- struct{}{}
	p.isMiningMutex.Unlock()
	return nil
}

func (p *PoWPeer) createBlock(verificationKey string, slot int, draw lottery_strategy.WinningLotteryParams, blocktree PoWblockchain.Blocktree) (newBlock PoWblockchain.Block, isEmpty bool) {
	//TODO Need to check that the draw is correct
	secretKey, foundSk := p.getSecretKey(verificationKey)
	if !foundSk {
		panic("Tried to create a block but peer did not have the associated SecretKey")
	}
	parentHash := draw.ParentHash
	unfinalizedTransactions := <-p.unfinalizedTransactions
	allTransactionsToAdd := blocktree.GetTransactionsNotInTree(unfinalizedTransactions)
	p.unfinalizedTransactions <- unfinalizedTransactions

	var transactionsToAdd []models.SignedTransaction
	if len(allTransactionsToAdd) <= p.maximumTransactionsInBlock {
		transactionsToAdd = allTransactionsToAdd
	}
	if len(allTransactionsToAdd) > p.maximumTransactionsInBlock {
		transactionsToAdd = make([]models.SignedTransaction, p.maximumTransactionsInBlock)
		for i := 0; i < p.maximumTransactionsInBlock; i++ {
			transactionsToAdd[i] = allTransactionsToAdd[i]
			//This could maybe cause starvation of transactions, if not enough blocks are made to saturate transaction demand
		}
	}
	//
	resultBlock := PoWblockchain.Block{
		IsGenesis: false,
		Vk:        verificationKey,
		Slot:      slot,
		Draw:      draw,
		BlockData: PoWblockchain.BlockData{
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

func (p *PoWPeer) sendBlockWithTransactions(draw lottery_strategy.WinningLotteryParams) {
	secretKeys := <-p.publicToSecret
	verificationKey := utils.GetSomeKey(secretKeys) //todo maybe make sure that it is the same public key that was used for the draw
	p.publicToSecret <- secretKeys
	blocktree := <-p.blockTreeChan
	extendedOnSlot := blocktree.HashToBlock(draw.ParentHash).Slot
	slot := utils.CalculateSlot(p.startTime)
	for slot <= extendedOnSlot {
		time.Sleep(constants.SlotLength / 10)
		slot = utils.CalculateSlot(p.startTime)
	}
	blockWithTransactions, isEmpty := p.createBlock(verificationKey, slot, draw, blocktree)
	if isEmpty {
		p.blockTreeChan <- blocktree
		return
	}
	msg := Message.Message{
		MessageType:      constants.BlockDelivery,
		MessageSender:    p.network.GetAddress().ToString(),
		PoWMessageBlocks: []PoWblockchain.Block{blockWithTransactions},
	}
	for _, block := range msg.PoWMessageBlocks {
		p.unhandledBlocks <- block
	}
	p.blockTreeChan <- blocktree
	p.network.FloodMessageToAllKnown(msg)
}

func (p *PoWPeer) blockHandlerLoop() {
	for {
		blockToHandle := <-p.unhandledBlocks
		go p.handleBlock(blockToHandle)
	}
}

func (p *PoWPeer) verifyBlock(block PoWblockchain.Block) bool {
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

func (p *PoWPeer) verifyTransactions(transactions []models.SignedTransaction) bool {
	for _, transaction := range transactions {
		transactionSignatureIsCorrect := utils.TransactionHasCorrectSignature(p.signatureStrategy, transaction)
		if !transactionSignatureIsCorrect {
			return false
		}
	}
	return true
}

func (p *PoWPeer) handleBlock(block PoWblockchain.Block) {
	if !p.verifyBlock(block) {
		return
	}
	blocktree := <-p.blockTreeChan
	block = Message.MakeDeepCopyOfPoWBlock(block)
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
		time.Sleep(1000 * time.Millisecond) //Needs to be enough time for the other block to arrive
		p.unhandledBlocks <- block
	case 1:
		//Block successfully added to the tree
		p.blockTreeChan <- blocktree
	default:
		p.blockTreeChan <- blocktree
		panic("addBlockReturnValueNotUnderstood")
	}
}

func (p *PoWPeer) addTransaction(t models.SignedTransaction) {
	unfinalizedTransactions := <-p.unfinalizedTransactions
	unfinalizedTransactions = append(unfinalizedTransactions, t)
	p.unfinalizedTransactions <- unfinalizedTransactions
}

func (p *PoWPeer) blockCreatingLoop(wins chan lottery_strategy.WinningLotteryParams) {
	for {
		newWin := <-wins
		go p.sendBlockWithTransactions(newWin)
	}
}
