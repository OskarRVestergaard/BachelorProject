package Peer

import (
	"fmt"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"github.com/OskarRVestergaard/BachelorProject/production/utils/constants"
	"math/big"
	"time"
)

// TODO FIX LATER
func (p *Peer) Mine() {
	fmt.Println("We're there. All I can see are turtle tracks. Whaddaya say we give Bowser the old Brooklyn one-two?")
	var hasPotentialWinner bool
	for k := range p.PublicToSecret {
		hasPotentialWinner, _ = p.lotteryStrategy.Mine(k, "PrevHash")
	}

	if hasPotentialWinner {
		// do block stuff IDK
	}
}

func (p *Peer) SendFakeBlockWithTransactions() {
	var publicKey = utils.GetSomeKey(p.PublicToSecret)
	var secretKey = p.PublicToSecret[publicKey]
	var headBlock = p.blockTree.GetHead()
	var headBlockHash = headBlock.HashOfBlock()
	var blockWithCurrentlyUnhandledTransactions = blockchain.Block{
		IsGenesis: false,
		Vk:        publicKey,
		Slot:      1,
		Draw:      "+4 cards",
		BlockData: blockchain.BlockData{
			Transactions: p.unfinalizedTransactions, //TODO Should only add not already added transactions (ones not in the chain) This is both something the create of the block should take care of, but also something that the receiver needs to check
		},
		ParentHash: headBlockHash,
		Signature:  big.Int{},
	}
	blockWithCurrentlyUnhandledTransactions.CalculateSignature(p.signatureStrategy, secretKey)
	var msg = blockchain.Message{
		MessageType:   constants.BlockDelivery,
		MessageSender: p.IpPort,
		MessageBlocks: []blockchain.Block{blockWithCurrentlyUnhandledTransactions},
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
	hasCorrectSignature := block.HasCorrectSignature(p.signatureStrategy)
	if !hasCorrectSignature {
		return
	}
	//Check correctness of transactions
	transactions := block.BlockData.Transactions
	for _, transaction := range transactions {
		if !utils.TransactionHasCorrectSignature(p.signatureStrategy, transaction) {
			return
		}
	}

	//ONLY ADD TRANSACTION THAT DO NOT ALREADY EXIST

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
